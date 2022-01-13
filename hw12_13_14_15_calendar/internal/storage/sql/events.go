package sqlstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
)

type SQLEventExpression struct {
	storage.BasicEventExpression
	db *sqlx.DB
}

func (ee SQLEventExpression) Execute(ctx context.Context, page int) *storage.EventIterator {
	return nil
}

type SQLEventRepository struct {
	db *sqlx.DB
}

func (S SQLEventRepository) Attach(db *sqlx.DB) {
	S.db = db
}

func (S SQLEventRepository) All(ctx context.Context, buffer []storage.Event) {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) One(ctx context.Context, id storage.ID) storage.Event {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Create(ctx context.Context, e storage.Event) (added storage.Event, err error) {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Update(ctx context.Context, e storage.Event) error {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Delete(ctx context.Context, e storage.Event) error {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Where() storage.EventExpression {
	return SQLEventExpression{
		db: S.db,
	}
}
