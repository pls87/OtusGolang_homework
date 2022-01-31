package grpc

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type eventStep struct {
	action      string
	e           models.Event
	expectedRes []models.Event
	expectedErr error
}

type bridgeTestSuite struct {
	suite.Suite
	bridge  *EventService
	storage basic.Storage
	now     time.Time
}

func (s *bridgeTestSuite) SetupTest() {
	l := &logrus.Logger{}
	s.storage = memory.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
	s.bridge = &EventService{
		logger:   l,
		eventApp: app.NewEventApp(s.storage, l),
	}
	s.now = time.Now().In(time.UTC)
}

func (s *bridgeTestSuite) TearDownTest() {
	s.NoError(s.storage.Dispose())
}

func (s *bridgeTestSuite) TestBasicOperations() {
	s.RunSteps(basicSteps)
}

func (s *bridgeTestSuite) TestFilterMonthOperation() {
	steps := seedSteps(s.now)
	s.RunSteps(steps)

	ec, err := s.bridge.GetEvents(context.Background(), &generated.Period{Unit: "month"})
	s.NoError(err)

	events := protoCollection2Events(ec)
	s.compareSlices(steps[2].expectedRes, events)
}

func (s *bridgeTestSuite) TestFilterWeekOperation() {
	steps := seedSteps(s.now)
	s.RunSteps(steps)

	ec, err := s.bridge.GetEvents(context.Background(), &generated.Period{Unit: "week"})
	s.NoError(err)

	events := protoCollection2Events(ec)
	s.compareSlices(steps[2].expectedRes[0:2], events)
}

func (s *bridgeTestSuite) TestFilterDayOperation() {
	steps := seedSteps(s.now)
	s.RunSteps(steps)

	ec, err := s.bridge.GetEvents(context.Background(), &generated.Period{Unit: "day"})
	s.NoError(err)

	events := protoCollection2Events(ec)
	l := 1
	if s.now.Weekday() == time.Monday {
		l = 2
	}
	s.compareSlices(steps[2].expectedRes[0:l], events)
}

func (s *bridgeTestSuite) RunSteps(steps []eventStep) {
	for _, st := range steps {
		switch st.action {
		case "create":
			added, err := s.bridge.AddEvent(context.Background(), event2Proto(st.e))
			if st.expectedErr != nil {
				s.ErrorIs(err, st.expectedErr)
				continue
			} else {
				s.NoError(err)
			}
			s.Greaterf(added.Id, int64(0), "ID of the added element should be more than 0")
			// check that event was added to storage correctly
			item, e := s.storage.Events().One(context.Background(), models.ID(added.Id))
			s.NoError(e)
			st.e.ID = models.ID(added.Id)
			s.Equal(st.e, item)
		case "update":
			updated, err := s.bridge.UpdateEvent(context.Background(), event2Proto(st.e))
			if st.expectedErr != nil {
				s.ErrorIs(err, st.expectedErr)
				continue
			} else {
				s.NoError(err)
			}
			// check that event was updated in storage correctly
			item, e := s.storage.Events().One(context.Background(), models.ID(updated.Id))
			s.NoError(e)
			s.Equal(st.e, item)
		case "delete":
			_, err := s.bridge.Delete(context.Background(), event2Proto(st.e))
			if st.expectedErr != nil {
				s.ErrorIs(err, st.expectedErr)
				continue
			} else {
				s.NoError(err)
			}
			// check that event was correctly removed from storage
			_, e := s.storage.Events().One(context.Background(), st.e.ID)
			s.ErrorIs(e, basic.ErrDoesNotExist)
		}
		fromBridge, err := s.bridge.GetAllEvents(context.Background(), &generated.Empty{})
		s.NoError(err)
		eventsFromStorage, e := s.storage.Events().All(context.Background())
		s.NoError(e)
		eventsFromBridge := protoCollection2Events(fromBridge)

		sort.Slice(eventsFromStorage, func(i, j int) bool { return eventsFromStorage[i].ID < eventsFromStorage[j].ID })
		sort.Slice(eventsFromBridge, func(i, j int) bool { return eventsFromBridge[i].ID < eventsFromBridge[j].ID })
		sort.Slice(st.expectedRes, func(i, j int) bool { return st.expectedRes[i].ID < st.expectedRes[j].ID })

		s.Equal(st.expectedRes, eventsFromStorage)
		s.Equal(st.expectedRes, eventsFromBridge)
	}
}

func (s *bridgeTestSuite) compareSlices(expected, actual []models.Event) {
	s.Equal(len(expected), len(actual))
	sort.Slice(expected, func(i, j int) bool { return expected[i].ID < expected[j].ID })
	sort.Slice(actual, func(i, j int) bool { return actual[i].ID < actual[j].ID })
	s.Equal(expected, actual)
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(bridgeTestSuite))
}
