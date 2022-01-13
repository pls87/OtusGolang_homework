package storage

import (
	"context"
	"time"
)

const MaxDuration = time.Second*1<<63 - 1

type EventIterator interface {
	Next() error
	Finished() bool
	Buffer() []Event
}

type EventRepository interface {
	All(ctx context.Context, buffer []Event)
	Where() EventExpression
	One(ctx context.Context, id ID) Event
	Create(ctx context.Context, e Event) (added Event, err error)
	Update(ctx context.Context, e Event) error
	Delete(ctx context.Context, e Event) error
}

type Timeframe struct {
	Start    time.Time
	Duration time.Duration
}

type Event struct {
	Timeframe
	ID           ID
	Title        string
	UserID       ID
	NotifyBefore time.Duration
	Desc         string
}

type BasicEventExpression struct {
	userID       ID
	starts       Timeframe
	intersection Timeframe
}

type EventExpression interface {
	User(id ID) EventExpression
	StartsIn(tf Timeframe) EventExpression
	StartsLater(d time.Time) EventExpression
	StartsBefore(d time.Time) EventExpression
	StartsLast(d time.Duration) EventExpression
	Intersects(tf Timeframe) EventExpression
	Execute(ctx context.Context, page int) *EventIterator
}

func (ee BasicEventExpression) Execute(_ context.Context, _ int) *EventIterator {
	panic("Abstract method called")
}

func (ee BasicEventExpression) User(id ID) EventExpression {
	ee.userID = id
	return &ee
}

func (ee BasicEventExpression) StartsIn(tf Timeframe) EventExpression {
	ee.starts = tf
	return &ee
}

func (ee BasicEventExpression) StartsLater(d time.Time) EventExpression {
	return ee.StartsIn(Timeframe{
		Start:    d,
		Duration: MaxDuration,
	})
}

func (ee BasicEventExpression) StartsBefore(d time.Time) EventExpression {
	minDate := time.Unix(0, 0)
	return ee.StartsIn(Timeframe{
		Start:    minDate,
		Duration: d.Sub(minDate),
	})
}

func (ee BasicEventExpression) StartsLast(d time.Duration) EventExpression {
	return ee.StartsIn(Timeframe{
		Start:    time.Now().Add(-d),
		Duration: d,
	})
}

func (ee BasicEventExpression) Intersects(tf Timeframe) EventExpression {
	ee.intersection = tf

	return &ee
}
