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

func (ee *BasicEventExpression) Execute(_ context.Context) (EventIterator, error) {
	panic("Abstract method called")
}

func (ee *BasicEventExpression) User(id models.ID) EventExpression {
	ee.UserID = id
	return ee
}

func (ee *BasicEventExpression) StartsIn(tf models.Timeframe) EventExpression {
	ee.Starts = tf
	return ee
}

func (ee *BasicEventExpression) StartsLater(d time.Time) EventExpression {
	return ee.StartsIn(models.Timeframe{
		Start:    d,
		Duration: models.MaxDuration,
	})
}

func (ee *BasicEventExpression) StartsBefore(d time.Time) EventExpression {
	minDate := time.Unix(0, 0)
	return ee.StartsIn(models.Timeframe{
		Start:    minDate,
		Duration: d.Sub(minDate),
	})
}

func (ee *BasicEventExpression) StartsLast(d time.Duration) EventExpression {
	return ee.StartsIn(models.Timeframe{
		Start:    time.Now().Add(-d),
		Duration: d,
	})
}

func (ee *BasicEventExpression) Intersects(tf models.Timeframe) EventExpression {
	ee.Intersection = tf

	return ee
}
