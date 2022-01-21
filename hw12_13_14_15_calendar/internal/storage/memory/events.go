package memorystorage

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type MemoryEventIterator struct {
	index int
	items []models.Event
	mu    *sync.RWMutex
}

func (s *MemoryEventIterator) Next() bool {
	s.index++
	return s.index < len(s.items)
}

func (s *MemoryEventIterator) Current() (models.Event, error) {
	if s.index < len(s.items) {
		return s.items[s.index], nil
	}
	return models.Event{}, fmt.Errorf("iterator is completed: %w", io.EOF)
}

func (s *MemoryEventIterator) ToArray() ([]models.Event, error) {
	return s.items, nil
}

func (s *MemoryEventIterator) Complete() error {
	return nil
}

type MemoryEventExpression struct {
	params *abstractstorage.EventExpressionParams
	mu     *sync.RWMutex
	data   *map[models.ID]models.Event
}

func (ee *MemoryEventExpression) Execute(_ context.Context) (abstractstorage.EventIterator, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	events := make([]models.Event, 0, 10)
	for _, v := range *ee.data {
		if ee.params.CheckEvent(v) {
			events = append(events, v)
		}
	}
	return &MemoryEventIterator{
		mu:    ee.mu,
		items: events,
	}, nil
}

func (ee *MemoryEventExpression) User(id models.ID) abstractstorage.EventExpression {
	ee.params.User(id)
	return ee
}

func (ee *MemoryEventExpression) StartsIn(tf models.Timeframe) abstractstorage.EventExpression {
	ee.params.StartsIn(tf)
	return ee
}

func (ee *MemoryEventExpression) StartsLater(d time.Time) abstractstorage.EventExpression {
	ee.params.StartsLater(d)
	return ee
}

func (ee *MemoryEventExpression) StartsBefore(d time.Time) abstractstorage.EventExpression {
	ee.params.StartsBefore(d)
	return ee
}

func (ee *MemoryEventExpression) StartsLast(d time.Duration) abstractstorage.EventExpression {
	ee.params.StartsLast(d)
	return ee
}

func (ee *MemoryEventExpression) Intersects(tf models.Timeframe) abstractstorage.EventExpression {
	ee.params.Intersects(tf)
	return ee
}

type MemoryEventRepository struct {
	mu      *sync.RWMutex
	data    map[models.ID]models.Event
	idIndex models.ID
}

func (ee *MemoryEventRepository) Init() {
	ee.data = make(map[models.ID]models.Event)
}

func (ee *MemoryEventRepository) All(_ context.Context) (abstractstorage.EventIterator, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	events := make([]models.Event, 0, len(ee.data))
	for _, v := range ee.data {
		events = append(events, v)
	}
	return &MemoryEventIterator{
		mu:    ee.mu,
		items: events,
	}, nil
}

func (ee *MemoryEventRepository) One(_ context.Context, id models.ID) (models.Event, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if val, ok := ee.data[id]; ok {
		return val, nil
	}

	return models.Event{}, fmt.Errorf("GET: event id=%d: %w", id, abstractstorage.ErrDoesNotExist)
}

func (ee *MemoryEventRepository) Create(_ context.Context, e models.Event) (added models.Event, err error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	e.ID = ee.idIndex + 1
	ee.data[ee.idIndex+1] = e
	ee.idIndex++
	return e, nil
}

func (ee *MemoryEventRepository) Update(_ context.Context, e models.Event) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.data[e.ID]; ok {
		ee.data[e.ID] = e
		return nil
	}

	return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, abstractstorage.ErrDoesNotExist)
}

func (ee *MemoryEventRepository) Delete(_ context.Context, e models.Event) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.data[e.ID]; ok {
		delete(ee.data, e.ID)
		return nil
	}

	return fmt.Errorf("DELETE: event id=%d: %w", e.ID, abstractstorage.ErrDoesNotExist)
}

func (ee *MemoryEventRepository) Where() abstractstorage.EventExpression {
	res := MemoryEventExpression{
		mu:     ee.mu,
		data:   &ee.data,
		params: &abstractstorage.EventExpressionParams{},
	}

	return &res
}
