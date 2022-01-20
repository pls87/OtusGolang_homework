package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type SQLEventExpression struct {
	abstractstorage.BasicEventExpression
	db *sqlx.DB
}

type SQLEventIterator struct {
	rows *sqlx.Rows
}

func (s *SQLEventIterator) Next() bool {
	return s.rows.Next()
}

func (s *SQLEventIterator) Current() (models.Event, error) {
	var ev models.Event
	e := s.rows.StructScan(&ev)
	return ev, e
}

func (s *SQLEventIterator) Complete() error {
	return s.rows.Close()
}

func (s *SQLEventIterator) ToArray() ([]models.Event, error) {
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

func (ee *SQLEventExpression) Execute(ctx context.Context) (abstractstorage.EventIterator, error) {
	clauseBuilder := make([]string, 0, 3)
	if ee.UserID > 0 {
		clauseBuilder = append(clauseBuilder, fmt.Sprintf("(user_id=%d)", ee.UserID))
	}
	if !ee.Starts.Start.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf("(start>=%s AND start<=%s)", ee.Starts.Start, ee.Starts.End()),
		)
	}

	if !ee.Intersection.Start.IsZero() {
		clauseBuilder = append(clauseBuilder,
			fmt.Sprintf("((start>=%s AND start<=%s) OR (start + duration >= %s AND start + duration <= %s))",
				ee.Intersection.Start, ee.Intersection.End(),
				ee.Intersection.Start, ee.Intersection.End(),
			),
		)
	}

	whereClause := strings.Join(clauseBuilder, " AND ")

	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	rows, err := ee.db.QueryxContext(ctx, `SELECT * FROM "events" ?`, whereClause) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &SQLEventIterator{rows}, nil
}

type SQLEventRepository struct {
	db *sqlx.DB
}

func (s *SQLEventRepository) Attach(db *sqlx.DB) {
	s.db = db
}

func (s *SQLEventRepository) Init() {
}

func (s *SQLEventRepository) All(ctx context.Context) (abstractstorage.EventIterator, error) {
	rows, err := s.db.QueryxContext(ctx, `SELECT * FROM "events"`) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &SQLEventIterator{rows}, nil
}

func (s *SQLEventRepository) One(ctx context.Context, id models.ID) (models.Event, error) {
	var ev models.Event
	err := s.db.GetContext(ctx, &ev, `SELECT * FROM events WHERE id=%d`, id)

	switch {
	case err == nil:
		return ev, nil
	case errors.Is(err, sql.ErrNoRows):
		return ev, fmt.Errorf("SELECT: event id=%d: %w", id, abstractstorage.ErrDoesNotExist)
	default:
		return ev, err
	}
}

func (s *SQLEventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
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

func (s *SQLEventRepository) Update(ctx context.Context, e models.Event) error {
	query := `UPDATE "events" SET 
        title='?', user_id=?, start=TIMESTAMP WITH TIME ZONE '?', 
        duration='? nanoseconds', notify_before='? nanoseconds',  description='?' WHERE id=?`
	res, err := s.db.ExecContext(ctx, query, e.Title, e.UserID, e.Start,
		e.Duration.Nanoseconds(), e.NotifyBefore.Nanoseconds(), e.Desc, e.ID)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, abstractstorage.ErrDoesNotExist)
		}
	}
	return err
}

func (s *SQLEventRepository) Delete(ctx context.Context, e models.Event) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM "events" WHERE id=?`, e.ID)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("DELETE: event id=%d: %w", e.ID, abstractstorage.ErrDoesNotExist)
		}
	}
	return err
}

func (s *SQLEventRepository) Where() abstractstorage.EventExpression {
	res := SQLEventExpression{
		db: s.db,
	}

	return &res
}
