package app

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/sirupsen/logrus"
)

type Application interface {
	Events() EventApplication
}

type App struct {
	events *EventApp
}

func (a *App) Events() EventApplication {
	return a.events
}

func New(logger *logrus.Logger, storage basic.Storage) Application {
	return &App{
		events: NewEventApp(storage, logger),
	}
}
