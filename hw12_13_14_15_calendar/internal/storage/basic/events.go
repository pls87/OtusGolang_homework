package basic

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
	All(ctx context.Context) ([]models.Event, error)
	Select() EventExpression
	One(ctx context.Context, id models.ID) (models.Event, error)
	Create(ctx context.Context, e models.Event) (added models.Event, err error)
	Update(ctx context.Context, e models.Event) error
	Delete(ctx context.Context, id models.ID) error
	DeleteObsolete(ctx context.Context, ttl time.Duration) error
}

type EventExpression interface {
	User(id models.ID) EventExpression
	StartsIn(tf models.Timeframe) EventExpression
	Intersects(tf models.Timeframe) EventExpression
	ToNotify() EventExpression
	Execute(ctx context.Context) (EventIterator, error)
}

type EventExpressionParams struct {
	UserID       models.ID
	Starts       models.Timeframe
	Intersection models.Timeframe
	ToNotify     time.Time
}

func (ee *EventExpressionParams) User(id models.ID) {
	ee.UserID = id
}

func (ee *EventExpressionParams) Notify() {
	ee.ToNotify = time.Now()
}

func (ee *EventExpressionParams) StartsIn(tf models.Timeframe) {
	ee.Starts = tf
}

func (ee *EventExpressionParams) Intersects(tf models.Timeframe) {
	ee.Intersection = tf
}

func (ee *EventExpressionParams) CheckEvent(e models.Event) bool {
	if ee.UserID > 0 && e.UserID != ee.UserID {
		return false
	}
	if !ee.ToNotify.IsZero() {
		return e.Start.After(ee.ToNotify) && e.Start.Sub(ee.ToNotify) < e.NotifyBefore
	}
	if !ee.Starts.Start.IsZero() && !(e.Start.After(ee.Starts.Start) && e.Start.Before(ee.Starts.End())) {
		return false
	}

	if !ee.Intersection.Start.IsZero() &&
		!((e.Start.After(ee.Intersection.Start) && e.Start.Before(ee.Intersection.End())) ||
			(e.Timeframe.End().After(ee.Intersection.Start) && e.Timeframe.End().Before(ee.Intersection.End()))) {
		return false
	}

	return true
}
