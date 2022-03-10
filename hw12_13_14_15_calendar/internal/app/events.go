package app

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventApplication struct {
	logger  *logrus.Logger
	storage basic.Storage
}

func NewEventApp(s basic.Storage, l *logrus.Logger) EventApplication {
	return EventApplication{
		logger:  l,
		storage: s,
	}
}

func (a *EventApplication) All(ctx context.Context) (collection []models.Event, err error) {
	return a.storage.Events().All(ctx)
}

func (a *EventApplication) New(ctx context.Context, e models.Event) (created models.Event, err error) {
	return a.storage.Events().Create(ctx, e)
}

func (a *EventApplication) Update(ctx context.Context, e models.Event) (err error) {
	return a.storage.Events().Update(ctx, e)
}

func (a *EventApplication) Remove(ctx context.Context, id models.ID) (err error) {
	return a.storage.Events().Delete(ctx, id)
}

func (a *EventApplication) ForTimeframe(ctx context.Context, f models.Timeframe) (events []models.Event, err error) {
	it, err := a.storage.Events().Select().Intersects(f).Execute(ctx)
	if err != nil {
		return nil, err
	}

	return it.ToArray()
}
