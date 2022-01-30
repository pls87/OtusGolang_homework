package handler

import (
	"net/http"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	resp   *response
	events *EventService
}

func NewService(app app.Application, logger *logrus.Logger) *Service {
	resp := &response{logger: logger}
	return &Service{
		events: &EventService{logger: logger, eventApp: app.Events(), resp: resp},
		logger: logger,
		resp:   resp,
	}
}

func (s *Service) Events() *EventService {
	return s.events
}

func (s *Service) Noop(w http.ResponseWriter, r *http.Request) {
	s.resp.text(r.Context(), w, "It Works!")
}
