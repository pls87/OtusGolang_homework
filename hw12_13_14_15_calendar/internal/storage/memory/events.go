package memory

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventIterator struct {
	index int
	items []models.Event
	mu    *sync.RWMutex
}

func (s *EventIterator) Next() bool {
	s.index++
	return s.index < len(s.items)
}

func (s *EventIterator) Current() (models.Event, error) {
	if s.index < len(s.items) {
		return s.items[s.index], nil
	}
	return models.Event{}, fmt.Errorf("iterator is completed: %w", io.EOF)
}

func (s *EventIterator) ToArray() ([]models.Event, error) {
	return s.items, nil
}

func (s *EventIterator) Complete() error {
	return nil
}

type EventExpression struct {
	params *basic.EventExpressionParams
	mu     *sync.RWMutex
	data   *map[models.ID]models.Event
}

func (ee *EventExpression) Execute(_ context.Context) (basic.EventIterator, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	events := make([]models.Event, 0, 10)
	for _, v := range *ee.data {
		if ee.params.CheckEvent(v) {
			events = append(events, v)
		}
	}
	return &EventIterator{
		mu:    ee.mu,
		items: events,
	}, nil
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

type EventRepository struct {
	mu      *sync.RWMutex
	data    map[models.ID]models.Event
	idIndex models.ID
}

func (ee *EventRepository) Init() {
	ee.data = make(map[models.ID]models.Event)
}

func (ee *EventRepository) All(_ context.Context) (events []models.Event, e error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	events = make([]models.Event, 0, len(ee.data))
	for _, v := range ee.data {
		events = append(events, v)
	}
	return events, nil
}

func (ee *EventRepository) One(_ context.Context, id models.ID) (models.Event, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if val, ok := ee.data[id]; ok {
		return val, nil
	}

	return models.Event{}, fmt.Errorf("GET: event id=%d: %w", id, basic.ErrDoesNotExist)
}

func (ee *EventRepository) Create(_ context.Context, e models.Event) (added models.Event, err error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	e.ID = ee.idIndex + 1
	ee.data[ee.idIndex+1] = e
	ee.idIndex++
	return e, nil
}

func (ee *EventRepository) Update(_ context.Context, e models.Event) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.data[e.ID]; ok {
		ee.data[e.ID] = e
		return nil
	}

	return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
}

func (ee *EventRepository) Delete(_ context.Context, e models.Event) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.data[e.ID]; ok {
		delete(ee.data, e.ID)
		return nil
	}

	return fmt.Errorf("DELETE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
}

func (ee *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		mu:     ee.mu,
		data:   &ee.data,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
