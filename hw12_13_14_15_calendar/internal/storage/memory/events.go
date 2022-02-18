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
	data   map[models.ID]models.Event
	sent   map[models.ID]bool
}

func (ee *EventExpression) checkEvent(e models.Event) bool {
	p := ee.params
	if p.UserID > 0 && e.UserID != p.UserID {
		return false
	}
	if !p.ToNotify.IsZero() &&
		((ee.sent[e.ID]) || (e.Start.Before(p.ToNotify) || e.Start.Sub(p.ToNotify) > e.NotifyBefore)) {
		return false
	}
	if !p.Starts.Start.IsZero() && !(e.Start.After(p.Starts.Start) && e.Start.Before(p.Starts.End())) {
		return false
	}

	if !p.Intersection.Start.IsZero() &&
		!((e.Start.After(p.Intersection.Start) && e.Start.Before(p.Intersection.End())) ||
			(e.Timeframe.End().After(p.Intersection.Start) && e.Timeframe.End().Before(p.Intersection.End()))) {
		return false
	}

	return true
}

func (ee *EventExpression) Execute(_ context.Context) (basic.EventIterator, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	events := make([]models.Event, 0, 10)
	for _, v := range ee.data {
		if ee.checkEvent(v) {
			events = append(events, v)
		}
	}
	return &EventIterator{
		index: -1,
		mu:    ee.mu,
		items: events,
	}, nil
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

type EventRepository struct {
	mu      *sync.RWMutex
	data    map[models.ID]models.Event
	sent    map[models.ID]bool
	idIndex models.ID
}

func (ee *EventRepository) Init() {
	ee.data = make(map[models.ID]models.Event)
	ee.sent = make(map[models.ID]bool)
}

func (ee *EventRepository) TrackSent(_ context.Context, eventID models.ID) (err error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.sent[eventID]; ok {
		return basic.ErrNotificationAlreadySent
	}
	ee.sent[eventID] = true
	return nil
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
	ee.idIndex++
	e.ID = ee.idIndex
	ee.data[ee.idIndex] = e
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

func (ee *EventRepository) Delete(_ context.Context, id models.ID) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	if _, ok := ee.data[id]; ok {
		delete(ee.data, id)
		return nil
	}

	return fmt.Errorf("DELETE: event id=%d: %w", id, basic.ErrDoesNotExist)
}

func (ee *EventRepository) DeleteObsolete(_ context.Context, ttl time.Duration) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	for k, v := range ee.data {
		if v.Start.Add(ttl).Before(time.Now()) {
			delete(ee.data, k)
		}
	}

	return nil
}

func (ee *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		mu:     ee.mu,
		data:   ee.data,
		sent:   ee.sent,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
