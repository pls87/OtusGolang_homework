package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
)

var ErrIncorrectParameter = errors.New("incorrect parameter received")

type EventService struct {
	generated.UnimplementedCalendarServer
	logger   *logrus.Logger
	eventApp app.EventApplication
}

func NewService(app app.EventApplication, logger *logrus.Logger) *EventService {
	return &EventService{
		logger:   logger,
		eventApp: app,
	}
}

func (es EventService) GetAllEvents(ctx context.Context, _ *generated.Empty) (*generated.EventCollection, error) {
	events, err := es.eventApp.All(ctx)
	if err != nil {
		return nil, err
	}

	return events2ProtoCollection(events), nil
}

func (es EventService) GetEvents(ctx context.Context, p *generated.Period) (*generated.EventCollection, error) {
	var frame models.Timeframe
	ok := frame.Period(time.Now(), p.GetUnit())
	if !ok {
		return nil, fmt.Errorf("%w: unknown period, awaiting 'day', 'week' or 'month'", ErrIncorrectParameter)
	}
	events, err := es.eventApp.ForTimeframe(ctx, frame)
	if err != nil {
		return nil, err
	}

	return events2ProtoCollection(events), nil
}

func (es EventService) AddEvent(ctx context.Context, pe *generated.Event) (*generated.Event, error) {
	e := proto2Event(pe)
	created, err := es.eventApp.New(ctx, e)
	if err != nil {
		return nil, err
	}

	return event2Proto(created), nil
}

func (es EventService) UpdateEvent(ctx context.Context, pe *generated.Event) (*generated.Event, error) {
	err := es.eventApp.Update(ctx, proto2Event(pe))
	if err != nil {
		return nil, err
	}

	return pe, nil
}

func (es EventService) Delete(ctx context.Context, pe *generated.Event) (*generated.Empty, error) {
	err := es.eventApp.Remove(ctx, models.ID(pe.Id))
	if err != nil {
		return nil, err
	}
	return &generated.Empty{}, nil
}
