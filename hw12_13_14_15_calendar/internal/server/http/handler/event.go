package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

type EventService struct {
	logger   *logrus.Logger
	eventApp app.EventApplication
	resp     *response
}

func (s *EventService) Get(w http.ResponseWriter, r *http.Request) {
	var frame *models.Timeframe
	var ok bool
	ctx := r.Context()
	if frame, ok = s.handleTimeframe(w, r); !ok {
		return
	}

	var events []models.Event
	var err error
	if frame == nil {
		events, err = s.eventApp.All(ctx)
	} else {
		events, err = s.eventApp.ForTimeframe(ctx, *frame)
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
	var id models.ID
	var status bool
	if id, status = s.handleID(w, r); !status {
		return
	}

	ctx := r.Context()

	toUpdate := models.Event{ID: id}
	err := json.NewDecoder(r.Body).Decode(&toUpdate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse event body", err)
		return
	}

	err = s.eventApp.Update(ctx, toUpdate)
	if err != nil {
		if errors.Is(err, basic.ErrDoesNotExist) {
			s.resp.notFound(ctx, w, fmt.Sprintf("event (id=%d) not found ", id), err)
			return
		}
		s.resp.internalServerError(ctx, w, "Unexpected error while saving event to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Event{"event": toUpdate})
}

func (s *EventService) Remove(w http.ResponseWriter, r *http.Request) {
	var id models.ID
	var status bool
	if id, status = s.handleID(w, r); !status {
		return
	}

	ctx := r.Context()
	err := s.eventApp.Remove(ctx, id)
	if err != nil {
		if errors.Is(err, basic.ErrDoesNotExist) {
			s.resp.notFound(ctx, w, fmt.Sprintf("event (id=%d) not found ", id), err)
			return
		}
		s.resp.internalServerError(ctx, w, "Unexpected error while saving event to storage", err)
		return
	}

	s.resp.json(ctx, w, true)
}

func (s *EventService) handleID(w http.ResponseWriter, r *http.Request) (models.ID, bool) {
	vars := mux.Vars(r)
	eventID, e := strconv.Atoi(vars["id"])
	if e != nil || eventID <= 0 {
		s.resp.badRequest(r.Context(), w, "malformed event id", e)
		return 0, false
	}

	return models.ID(eventID), true
}

func (s *EventService) handleTimeframe(w http.ResponseWriter, r *http.Request) (frame *models.Timeframe, ok bool) {
	frame = &models.Timeframe{}
	period := r.URL.Query()["period"]

	if len(period) > 0 {
		if res := frame.Period(time.Now(), period[0]); res {
			return frame, true
		}
		s.resp.badRequest(r.Context(), w, "Malformed parameter 'period'", ErrMalformedParam)
		return nil, false
	}

	return nil, true
}
