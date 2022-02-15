package sql

// haven't tested this package good enough

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
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

func (ee *EventExpression) ToNotify() basic.EventExpression {
	ee.params.Notify()
	return ee
}

func (ee *EventExpression) StartsIn(tf models.Timeframe) basic.EventExpression {
	ee.params.StartsIn(tf)
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
		if e != nil {
			_ = s.Complete()
			return nil, e
		}
	}

	return res, nil
}

// Execute TODO: clean up this code later.
func (ee *EventExpression) Execute(ctx context.Context) (basic.EventIterator, error) {
	clauseBuilder := make([]string, 0, 4)
	clauseArgs := make([]interface{}, 0, 9)
	ind := 0
	if ee.params.UserID > 0 {
		clauseBuilder = append(clauseBuilder, fmt.Sprintf("(user_id=$%d)", ind+1))
		clauseArgs = append(clauseArgs, ee.params.UserID)
		ind++
	}
	if !ee.params.ToNotify.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf(`(start>$%d AND start - notify_before < $%d 
				AND "ID" NOT IN (SELECT event_id from "notification_sent"))`, ind+1, ind+2))
		clauseArgs = append(clauseArgs, ee.params.ToNotify, ee.params.ToNotify)
		ind += 2
	}
	if !ee.params.Starts.Start.IsZero() {
		clauseBuilder = append(clauseBuilder, fmt.Sprintf("(start>=$%d AND start<=$%d)", ind+1, ind+2))
		clauseArgs = append(clauseArgs, ee.params.Starts.Start, ee.params.Starts.End())
		ind += 2
	}

	if !ee.params.Intersection.Start.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf("((start>=$%d AND start<=$%d) OR (start + duration >= $%d AND start + duration <= $%d))",
				ind+1, ind+2, ind+3, ind+4))
		clauseArgs = append(clauseArgs, ee.params.Intersection.Start, ee.params.Intersection.End(),
			ee.params.Intersection.Start, ee.params.Intersection.End())
	}

	whereClause := strings.Join(clauseBuilder, " AND ")

	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	rows, err := ee.db.QueryxContext(ctx, `SELECT * FROM "events"`+whereClause, clauseArgs...) //nolint:sqlclosecheck
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
	err := s.db.GetContext(ctx, &ev, `SELECT * FROM events WHERE ID=$1`, id)

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
                VALUES ($1, $2, TIMESTAMP WITH TIME ZONE $3, $4, $5, $6) RETURNING "ID"`
	lastID := 0
	err = s.db.QueryRowxContext(
		ctx, query, e.Title, e.UserID, e.Start, fmt.Sprintf("%d nanoseconds", e.Duration.Nanoseconds()),
		fmt.Sprintf("%d nanoseconds", e.NotifyBefore.Nanoseconds()), e.Desc).Scan(&lastID)

	if err == nil {
		e.ID = models.ID(lastID)
	}

	return e, err
}

func (s *EventRepository) Update(ctx context.Context, e models.Event) error {
	query := `UPDATE "events" SET 
        title=$1, user_id=?, start=TIMESTAMP WITH TIME ZONE $2, 
        duration=$3, notify_before=$4,  description=$5 WHERE ID=$6`
	res, err := s.db.ExecContext(ctx, query, e.Title, e.UserID, e.Start,
		fmt.Sprintf("%d nanoseconds", e.Duration.Nanoseconds()),
		fmt.Sprintf("%d nanoseconds", e.NotifyBefore.Nanoseconds()), e.Desc, e.ID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM "events" WHERE ID=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("DELETE: event id=%d: %w", id, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) DeleteObsolete(ctx context.Context, ttl time.Duration) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM "events" WHERE start + $1 < $2`,
		fmt.Sprintf("%d nanoseconds", ttl.Nanoseconds()), time.Now())
	return err
}

func (s *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		db:     s.db,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
