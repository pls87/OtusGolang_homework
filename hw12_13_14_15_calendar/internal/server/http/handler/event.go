package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventService struct {
	logger   *logrus.Logger
	eventApp app.EventApplication
	resp     *response
}

func (s *EventService) Get(w http.ResponseWriter, r *http.Request) {
	var events []models.Event
	var err error
	ctx := r.Context()
	if frame, e := s.handleTimeframe(r); e != nil {
		events, err = s.eventApp.All(ctx)
	} else {
		events, err = s.eventApp.ForTimeframe(ctx, frame)
	}

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

func (s *EventService) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toUpdate models.Event
	err := json.NewDecoder(r.Body).Decode(&toUpdate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse event body", err)
		return
	}

	err = s.eventApp.Update(ctx, toUpdate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving event to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Event{"event": toUpdate})
}

func (s *EventService) handleTimeframe(r *http.Request) (models.Timeframe, error) {
	frame := models.Timeframe{}
	now := time.Now()
	period := r.URL.Query()["period"]
	if len(period) > 0 {
		switch period[0] {
		case "day":
			frame.Start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			frame.Duration = 24 * time.Hour
		case "week":
			d := now
			for d.Weekday() != time.Monday {
				d = d.AddDate(0, 0, -1)
			}
			frame.Start = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
			frame.Duration = 7 * 24 * time.Hour
		case "month":
			frame.Start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
			frame.Start.AddDate(0, 1, 0)
			frame.Duration = frame.Start.AddDate(0, 1, 0).Sub(frame.Start)
		}
		return frame, nil
	}
	return frame, fmt.Errorf("%w: 'period'", ErrMissedRequiredParam)
}
