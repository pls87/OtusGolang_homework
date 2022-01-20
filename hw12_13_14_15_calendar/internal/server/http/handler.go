package internalhttp

import (
	"net/http"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/sirupsen/logrus"
)

type Service struct {
	app    app.Application
	logger *logrus.Logger
}

func NewService(app app.Application, logger *logrus.Logger) *Service {
	return &Service{
		app:    app,
		logger: logger,
	}
}

func (s *Service) Noop(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("It works!"))
}
