package abstractstorage

import (
	"context"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventIterator interface {
	Next() bool
	Current() (models.Event, error)
	ToArray() ([]models.Event, error)
	Complete() error
}

type EventRepository interface {
	Init()
	All(ctx context.Context) (EventIterator, error)
	Where() EventExpression
	One(ctx context.Context, id models.ID) (models.Event, error)
	Create(ctx context.Context, e models.Event) (added models.Event, err error)
	Update(ctx context.Context, e models.Event) error
	Delete(ctx context.Context, e models.Event) error
}

type BasicEventExpression struct {
	UserID       models.ID
	Starts       models.Timeframe
	Intersection models.Timeframe
}

type EventExpression interface {
	User(id models.ID) EventExpression
	StartsIn(tf models.Timeframe) EventExpression
	StartsLater(d time.Time) EventExpression
	StartsBefore(d time.Time) EventExpression
	StartsLast(d time.Duration) EventExpression
	Intersects(tf models.Timeframe) EventExpression
	Execute(ctx context.Context) (EventIterator, error)
}

type EventExpressionParams struct {
	UserID       models.ID
	Starts       models.Timeframe
	Intersection models.Timeframe
}

func (ee *EventExpressionParams) User(id models.ID) {
	ee.UserID = id
}

func (ee *EventExpressionParams) StartsIn(tf models.Timeframe) {
	ee.Starts = tf
}

func (ee *EventExpressionParams) StartsLater(d time.Time) {
	ee.StartsIn(models.Timeframe{
		Start:    d,
		Duration: models.MaxDuration,
	})
}

func (ee *EventExpressionParams) StartsBefore(d time.Time) {
	minDate := time.Unix(0, 0)
	ee.StartsIn(models.Timeframe{
		Start:    minDate,
		Duration: d.Sub(minDate),
	})
}

func (ee *EventExpressionParams) StartsLast(d time.Duration) {
	ee.StartsIn(models.Timeframe{
		Start:    time.Now().Add(-d),
		Duration: d,
	})
}

func (ee *EventExpressionParams) Intersects(tf models.Timeframe) {
	ee.Intersection = tf
}

func (ee *EventExpressionParams) CheckEvent(e models.Event) bool {
	if ee.UserID > 0 && e.UserID != ee.UserID {
		return false
	}
	if !ee.Starts.Start.IsZero() && !(e.Start.After(ee.Starts.Start) && e.Start.Before(ee.Starts.End())) {
		return false
	}

	if !ee.Intersection.Start.IsZero() && !((e.Start.After(ee.Starts.Start) && e.Start.Before(ee.Starts.End())) ||
		(e.Timeframe.End().After(ee.Intersection.Start) && e.Timeframe.End().Before(ee.Intersection.End()))) {
		return false
	}

	return true
}
