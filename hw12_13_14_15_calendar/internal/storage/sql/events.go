package sql

// haven't tested this package good enough

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventExpression struct {
	params *basic.EventExpressionParams
	db     *sqlx.DB
}

func (ee *EventExpression) User(id models.ID) basic.EventExpression {
	ee.params.User(id)
	return ee
}

func (ee *EventExpression) StartsIn(tf models.Timeframe) basic.EventExpression {
	ee.params.StartsIn(tf)
	return ee
}

func (ee *EventExpression) StartsLater(d time.Time) basic.EventExpression {
	ee.params.StartsLater(d)
	return ee
}

func (ee *EventExpression) StartsBefore(d time.Time) basic.EventExpression {
	ee.params.StartsBefore(d)
	return ee
}

func (ee *EventExpression) StartsLast(d time.Duration) basic.EventExpression {
	ee.params.StartsLast(d)
	return ee
}

func (ee *EventExpression) Intersects(tf models.Timeframe) basic.EventExpression {
	ee.params.Intersects(tf)
	return ee
}

type EventIterator struct {
	rows *sqlx.Rows
}

func (s *EventIterator) Next() bool {
	return s.rows.Next()
}

func (s *EventIterator) Current() (models.Event, error) {
	var ev models.Event
	e := s.rows.StructScan(&ev)
	return ev, e
}

func (s *EventIterator) Complete() error {
	return s.rows.Close()
}

func (s *EventIterator) ToArray() ([]models.Event, error) {
	res := make([]models.Event, 0, 10)
	var ev models.Event
	for s.rows.Next() {
		e := s.rows.StructScan(&ev)
		if e == nil {
			res = append(res, ev)
		}
	}

	return res, nil
}

func (ee *EventExpression) Execute(ctx context.Context) (basic.EventIterator, error) {
	clauseBuilder := make([]string, 0, 3)
	clauseArgs := make([]interface{}, 0, 5)
	if ee.params.UserID > 0 {
		clauseBuilder = append(clauseBuilder, "(user_id=?)")
		clauseArgs = append(clauseArgs, ee.params.UserID)
	}
	if !ee.params.Starts.Start.IsZero() {
		clauseBuilder = append(clauseBuilder, "(start>=? AND start<=?)")
		clauseArgs = append(clauseArgs, ee.params.Starts.Start, ee.params.Starts.End())
	}

	if !ee.params.Intersection.Start.IsZero() {
		clauseBuilder = append(clauseBuilder,
			"((start>=? AND start<=?) OR (start + duration >= ? AND start + duration <= ?))")
		clauseArgs = append(clauseArgs, ee.params.Intersection.Start, ee.params.Intersection.End(),
			ee.params.Intersection.Start, ee.params.Intersection.End())
	}

	whereClause := strings.Join(clauseBuilder, " AND ")

	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	rows, err := ee.db.QueryxContext(ctx, `SELECT * FROM "events"`+whereClause, clauseArgs) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &EventIterator{rows}, nil
}

type EventRepository struct {
	db *sqlx.DB
}

func (s *EventRepository) Init() {
}

func (s *EventRepository) All(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := s.db.SelectContext(ctx, &events, `SELECT * FROM "events"`)

	return events, err
}

func (s *EventRepository) One(ctx context.Context, id models.ID) (models.Event, error) {
	var ev models.Event
	err := s.db.GetContext(ctx, &ev, `SELECT * FROM events WHERE ID=%d`, id)

	switch {
	case err == nil:
		return ev, nil
	case errors.Is(err, sql.ErrNoRows):
		return ev, fmt.Errorf("SELECT: event id=%d: %w", id, basic.ErrDoesNotExist)
	default:
		return ev, err
	}
}

func (s *EventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
	query := `INSERT INTO "events" (title, user_id, start, duration, notify_before,  description) 
		VALUES ('?', ?, TIMESTAMP WITH TIME ZONE '?', '? nanoseconds', '? nanoseconds', '?')`
	res, err := s.db.ExecContext(
		ctx, query, e.Title, e.UserID, e.Start, e.Duration.Nanoseconds(), e.NotifyBefore.Nanoseconds(), e.Desc)
	if err == nil {
		id, _ := res.LastInsertId()
		e.ID = models.ID(id)
	}

	return e, err
}

func (s *EventRepository) Update(ctx context.Context, e models.Event) error {
	query := `UPDATE "events" SET 
        title='?', user_id=?, start=TIMESTAMP WITH TIME ZONE '?', 
        duration='? nanoseconds', notify_before='? nanoseconds',  description='?' WHERE ID=?`
	res, err := s.db.ExecContext(ctx, query, e.Title, e.UserID, e.Start,
		e.Duration.Nanoseconds(), e.NotifyBefore.Nanoseconds(), e.Desc, e.ID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) Delete(ctx context.Context, e models.Event) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM "events" WHERE ID=?`, e.ID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("DELETE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		db:     s.db,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
