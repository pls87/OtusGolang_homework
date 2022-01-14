package memorystorage

import (
	"context"
	"sync"

	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type MemoryEventExpression struct {
	abstractstorage.BasicEventExpression
	mu *sync.RWMutex
}

func (ee MemoryEventExpression) Execute(ctx context.Context, page int) *abstractstorage.EventIterator {
	return nil
}

type MemoryEventRepository struct {
	mu *sync.RWMutex
}

func (ee *MemoryEventRepository) All(ctx context.Context, buffer []models.Event) {
	// TODO implement me
	panic("implement me")
}

func (ee *MemoryEventRepository) One(ctx context.Context, id models.ID) models.Event {
	// TODO implement me
	panic("implement me")
}

func (ee *MemoryEventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
	// TODO implement me
	panic("implement me")
}

func (ee *MemoryEventRepository) Update(ctx context.Context, e models.Event) error {
	// TODO implement me
	panic("implement me")
}

func (ee *MemoryEventRepository) Delete(ctx context.Context, e models.Event) error {
	// TODO implement me
	panic("implement me")
}

func (ee *MemoryEventRepository) Where() abstractstorage.EventExpression {
	return MemoryEventExpression{
		mu: ee.mu,
	}
}
