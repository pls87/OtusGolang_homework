package events

import (
	"context"
	"sync"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

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
	if !p.ToNotify.IsZero() && (ee.sent[e.ID] || !p.Notify(e)) {
		return false
	}

	if !p.Starts.Start.IsZero() && !p.StartsIn(e) {
		return false
	}

	if !p.Intersection.Start.IsZero() && !p.Intersects(e) {
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
	ee.params.UserID = id
	return ee
}

func (ee *EventExpression) ToNotify() basic.EventExpression {
	ee.params.ToNotify = time.Now()
	return ee
}

func (ee *EventExpression) StartsIn(tf models.Timeframe) basic.EventExpression {
	ee.params.Starts = tf
	return ee
}

func (ee *EventExpression) Intersects(tf models.Timeframe) basic.EventExpression {
	ee.params.Intersection = tf
	return ee
}
