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

func (ee SQLEventExpression) Execute(ctx context.Context, page int) *abstractstorage.EventIterator {
	return nil
}

type SQLEventRepository struct {
	db *sqlx.DB
}

func (S SQLEventRepository) Attach(db *sqlx.DB) {
	S.db = db
}

func (S SQLEventRepository) All(ctx context.Context, buffer []models.Event) {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) One(ctx context.Context, id models.ID) models.Event {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Update(ctx context.Context, e models.Event) error {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Delete(ctx context.Context, e models.Event) error {
	// TODO implement me
	panic("implement me")
}

func (S SQLEventRepository) Where() abstractstorage.EventExpression {
	return SQLEventExpression{
		db: S.db,
	}
}
