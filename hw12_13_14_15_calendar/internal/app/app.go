package app

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
	"github.com/sirupsen/logrus"
)

type Application interface{}

type App struct {
	logger  *logrus.Logger
	storage abstractstorage.Storage
	cfg     config.Config
}

func New(logger *logrus.Logger, storage abstractstorage.Storage, cfg config.Config) *App {
	return &App{logger, storage, cfg}
}
