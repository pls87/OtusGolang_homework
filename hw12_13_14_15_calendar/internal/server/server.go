package server

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger, app app.Application, cfg configs.APIConf) basic.Server {
	var srv basic.Server
	switch cfg.Type {
	case "http":
		srv = http.New(logger, app, cfg)
	case "grpc":
		srv = grpc.New(logger, app, cfg)
	default:
		srv = http.New(logger, app, cfg)
	}

	return srv
}
