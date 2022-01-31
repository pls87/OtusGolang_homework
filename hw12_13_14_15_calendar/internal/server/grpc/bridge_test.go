package grpc

import (
	"context"
	"sort"
	"testing"

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
}

func (s *bridgeTestSuite) SetupTest() {
	l := &logrus.Logger{}
	s.storage = memory.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
	s.bridge = &EventService{
		logger:   l,
		eventApp: app.NewEventApp(s.storage, l),
	}
}

func (s *bridgeTestSuite) TearDownTest() {
	s.NoError(s.storage.Dispose())
}

func (s *bridgeTestSuite) TestBasicOperations() {
	s.RunSteps(basicSteps)
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

func TestStorage(t *testing.T) {
	suite.Run(t, new(bridgeTestSuite))
}
