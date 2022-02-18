package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventRepository struct {
	M       *sync.RWMutex
	data    map[models.ID]models.Event
	sent    map[models.ID]bool
	idIndex models.ID
}

func (ee *EventRepository) Init() {
	ee.data = make(map[models.ID]models.Event)
	ee.sent = make(map[models.ID]bool)
}

func (ee *EventRepository) TrackSent(_ context.Context, eventID models.ID) (err error) {
	ee.M.Lock()
	defer ee.M.Unlock()
	if _, ok := ee.sent[eventID]; ok {
		return basic.ErrNotificationAlreadySent
	}
	ee.sent[eventID] = true
	return nil
}

func (ee *EventRepository) All(_ context.Context) (events []models.Event, e error) {
	ee.M.Lock()
	defer ee.M.Unlock()
	events = make([]models.Event, 0, len(ee.data))
	for _, v := range ee.data {
		events = append(events, v)
	}
	return events, nil
}

func (ee *EventRepository) One(_ context.Context, id models.ID) (models.Event, error) {
	ee.M.Lock()
	defer ee.M.Unlock()

	if val, ok := ee.data[id]; ok {
		return val, nil
	}

	return models.Event{}, fmt.Errorf("GET: event id=%d: %w", id, basic.ErrDoesNotExist)
}

func (ee *EventRepository) Create(_ context.Context, e models.Event) (added models.Event, err error) {
	ee.M.Lock()
	defer ee.M.Unlock()
	ee.idIndex++
	e.ID = ee.idIndex
	ee.data[ee.idIndex] = e
	return e, nil
}

func (ee *EventRepository) Update(_ context.Context, e models.Event) error {
	ee.M.Lock()
	defer ee.M.Unlock()
	if _, ok := ee.data[e.ID]; ok {
		ee.data[e.ID] = e
		return nil
	}

	return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
}

func (ee *EventRepository) Delete(_ context.Context, id models.ID) error {
	ee.M.Lock()
	defer ee.M.Unlock()
	if _, ok := ee.data[id]; ok {
		delete(ee.data, id)
		return nil
	}

	return fmt.Errorf("DELETE: event id=%d: %w", id, basic.ErrDoesNotExist)
}

func (ee *EventRepository) DeleteObsolete(_ context.Context, ttl time.Duration) error {
	ee.M.Lock()
	defer ee.M.Unlock()
	for k, v := range ee.data {
		if v.Start.Add(ttl).Before(time.Now()) {
			delete(ee.data, k)
		}
	}

	return nil
}

func (ee *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		mu:     ee.M,
		data:   ee.data,
		sent:   ee.sent,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
