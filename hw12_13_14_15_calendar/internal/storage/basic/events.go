package basic

import (
	"context"
	"errors"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

var ErrNotificationAlreadySent = errors.New("notification already sent")

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
	TrackSent(ctx context.Context, ID models.ID) error
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

func (ee *EventExpressionParams) Intersects(e models.Event) bool {
	return e.Start.After(ee.Intersection.Start) && e.Start.Before(ee.Intersection.End()) ||
		ee.Intersection.Start.After(e.Start) && ee.Intersection.Start.Before(e.End())
}

func (ee *EventExpressionParams) Notify(e models.Event) bool {
	return e.Start.After(ee.ToNotify) && e.Start.Sub(ee.ToNotify) < e.NotifyBefore
}

func (ee *EventExpressionParams) StartsIn(e models.Event) bool {
	return e.Start.After(ee.Starts.Start) && e.Start.Before(ee.Starts.End())
}
