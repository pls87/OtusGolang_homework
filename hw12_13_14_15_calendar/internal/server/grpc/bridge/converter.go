package bridge

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func event2Proto(e models.Event) *generated.Event {
	return &generated.Event{
		Id:           int64(e.ID),
		Title:        e.Title,
		Desc:         e.Desc,
		Start:        timestamppb.New(e.Start),
		Duration:     durationpb.New(e.Duration),
		NotifyBefore: durationpb.New(e.NotifyBefore),
		UserId:       int64(e.UserID),
	}
}

func events2ProtoCollection(events []models.Event) *generated.EventCollection {
	res := &generated.EventCollection{
		Events: make([]*generated.Event, 0, 10),
	}

	for _, e := range events {
		res.Events = append(res.Events, event2Proto(e))
	}

	return res
}

func proto2Event(e *generated.Event) models.Event {
	return models.Event{
		ID:    models.ID(e.Id),
		Title: e.Title,
		Desc:  e.Desc,
		Timeframe: models.Timeframe{
			Start:    e.Start.AsTime(),
			Duration: e.Duration.AsDuration(),
		},
		NotifyBefore: e.NotifyBefore.AsDuration(),
		UserID:       models.ID(e.UserId),
	}
}

func protoCollection2Events(ec *generated.EventCollection) []models.Event {
	events := make([]models.Event, 0, 10)

	for _, e := range ec.Events {
		events = append(events, proto2Event(e))
	}

	return events
}
