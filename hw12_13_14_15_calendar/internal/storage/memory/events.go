package memorystorage

import (
	"context"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
)

type MemoryEventExpression struct {
	storage.BasicEventExpression
	mu sync.RWMutex
}

func (ee MemoryEventExpression) Execute(ctx context.Context, page int) *storage.EventIterator {
	return nil
}

type MemoryEventRepository struct {
	mu sync.RWMutex
}

func (ee MemoryEventRepository) All(ctx context.Context, buffer []storage.Event) {
	// TODO implement me
	panic("implement me")
}

func (ee MemoryEventRepository) One(ctx context.Context, id storage.ID) storage.Event {
	// TODO implement me
	panic("implement me")
}

func (ee MemoryEventRepository) Create(ctx context.Context, e storage.Event) (added storage.Event, err error) {
	// TODO implement me
	panic("implement me")
}

func (ee MemoryEventRepository) Update(ctx context.Context, e storage.Event) error {
	// TODO implement me
	panic("implement me")
}

func (ee MemoryEventRepository) Delete(ctx context.Context, e storage.Event) error {
	// TODO implement me
	panic("implement me")
}

func (ee MemoryEventRepository) Where() storage.EventExpression {
	return MemoryEventExpression{
		mu: ee.mu,
	}
}
