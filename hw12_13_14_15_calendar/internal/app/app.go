package app

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/sirupsen/logrus"
)

type Application interface{}

type App struct {
	logger  *logrus.Logger
	storage basicstorage.Storage
	cfg     configs.Config
}

func New(logger *logrus.Logger, storage basicstorage.Storage, cfg configs.Config) *App {
	return &App{logger, storage, cfg}
}
