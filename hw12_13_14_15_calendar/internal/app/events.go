package app

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventApplication interface {
	All(ctx context.Context) ([]models.Event, error)
	New(ctx context.Context, e models.Event) (created models.Event, err error)
	Update(ctx context.Context, e models.Event) error
	Remove(ctx context.Context, id models.ID) error
	ForTimeframe(ctx context.Context, timeframe models.Timeframe) (events []models.Event, err error)
}

type EventApp struct {
	logger  *logrus.Logger
	storage basic.Storage
}

func NewEventApp(s basic.Storage, l *logrus.Logger) *EventApp {
	return &EventApp{
		logger:  l,
		storage: s,
	}
}

func (a *EventApp) All(ctx context.Context) (collection []models.Event, err error) {
	return a.storage.Events().All(ctx)
}

func (a *EventApp) New(ctx context.Context, e models.Event) (created models.Event, err error) {
	return a.storage.Events().Create(ctx, e)
}

func (a *EventApp) Update(ctx context.Context, e models.Event) (err error) {
	return a.storage.Events().Update(ctx, e)
}

func (a *EventApp) Remove(ctx context.Context, id models.ID) (err error) {
	return a.storage.Events().Delete(ctx, id)
}

func (a *EventApp) ForTimeframe(ctx context.Context, frame models.Timeframe) (events []models.Event, err error) {
	it, err := a.storage.Events().Select().Intersects(frame).Execute(ctx)
	if err != nil {
		return nil, err
	}

	return it.ToArray()
}
