package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventService struct {
	logger   *logrus.Logger
	eventApp app.EventApplication
	resp     *response
}

func (s *EventService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := s.eventApp.All(ctx)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting events from storage", err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Event{"events": events})
}

func (s *EventService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Event
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse event body", err)
		return
	}

	created, err := s.eventApp.New(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving event to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Event{"event": created})
}
