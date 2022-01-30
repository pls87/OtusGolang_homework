package app

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventApplication interface {
	All(ctx context.Context) ([]models.Event, error)
	New(ctx context.Context, c models.Event) (created models.Event, err error)
}

type EventApp struct {
	logger  *logrus.Logger
	storage basic.Storage
}

func (a *EventApp) All(ctx context.Context) (collection []models.Event, err error) {
	return a.storage.Events().All(ctx)
}

func (a *EventApp) New(ctx context.Context, c models.Event) (created models.Event, err error) {
	return a.storage.Events().Create(ctx, c)
}
