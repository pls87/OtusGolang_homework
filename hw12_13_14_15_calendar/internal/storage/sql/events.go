package sqlstorage

import (
	"context"

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

func (ee *SQLEventExpression) Execute(ctx context.Context) abstractstorage.EventIterator {
	return nil
}

type SQLEventRepository struct {
	db *sqlx.DB
}

func (s *SQLEventRepository) Attach(db *sqlx.DB) {
	s.db = db
}

func (s *SQLEventRepository) All(ctx context.Context) (abstractstorage.EventIterator, error) {
	rows, err := s.db.QueryxContext(ctx, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}

	return &SQLEventIterator{rows}, nil
}

func (s *SQLEventRepository) One(ctx context.Context, id models.ID) (models.Event, error) {
	var ev models.Event
	row := s.db.QueryRowxContext(ctx, "SELECT * FROM events WHERE id=%d", id)

	if row.Err() != nil {
		return ev, row.Err()
	}

	e := row.StructScan(&ev)
	return ev, e
}

func (s *SQLEventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
	query := `INSERT INTO "events" (title, user_id, start, duration, notify_before,  description) 
		VALUES ('?', ?, TIMESTAMP WITH TIME ZONE '?', '? nanoseconds', '? nanoseconds', '?')`
	res, err := s.db.ExecContext(ctx, query, e.Title, e.UserID, e.Start, e.Duration.Nanoseconds(), e.NotifyBefore.Nanoseconds(), e.Desc)

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
	_, err := s.db.ExecContext(ctx, query, e.Title, e.UserID, e.Start,
		e.Duration.Nanoseconds(), e.NotifyBefore.Nanoseconds(), e.Desc, e.ID)
	return err
}

func (s *SQLEventRepository) Delete(ctx context.Context, e models.Event) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM "events" WHERE id=?`, e.ID)
	return err
}

func (s *SQLEventRepository) Where() abstractstorage.EventExpression {
	res := SQLEventExpression{
		db: s.db,
	}

	return &res
}
