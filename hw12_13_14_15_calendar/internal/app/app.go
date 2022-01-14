package app

import (
	"context"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"github.com/sirupsen/logrus"
)

type App struct { // TODO
	logger  *logrus.Logger
	storage abstractstorage.Storage
	cfg     config.Config
}

func New(logger *logrus.Logger, storage abstractstorage.Storage, cfg config.Config) *App {
	return &App{logger, storage, cfg}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
