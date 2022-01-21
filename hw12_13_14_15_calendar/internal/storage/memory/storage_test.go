package memorystorage_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/suite"
)

type eventStep struct {
	action      string
	e           models.Event
	expectedRes []models.Event
	expectedErr error
}

type eventsStorageTestSuite struct {
	suite.Suite
	storage basicstorage.Storage
}

func (s *eventsStorageTestSuite) SetupTest() {
	s.storage = memorystorage.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
}

func (s *eventsStorageTestSuite) TearDownTest() {
	s.NoError(s.storage.Dispose())
}

func (s *eventsStorageTestSuite) TestBasicOperations() {
	s.RunSteps(basicSteps)
}

func (s *eventsStorageTestSuite) TestFilterStartOperation() {
	s.RunSteps(seedSteps)

	exp := s.storage.Events().Select().StartsIn(models.Timeframe{
		Start:    time.Date(2022, 2, 1, 0, 0, 0, 0, time.Local),
		Duration: 30 * 24 * time.Hour,
	})

	s.ExecuteFilterResults(exp, seedSteps[1].expectedRes[1])
}

func (s *eventsStorageTestSuite) TestFilterIntersectionOperation() {
	s.RunSteps(seedSteps)

	// Test intersection start
	exp := s.storage.Events().Select().Intersects(models.Timeframe{
		Start:    time.Date(2022, 1, 12, 12, 30, 0, 0, time.Local),
		Duration: 45 * time.Minute,
	})
	s.ExecuteFilterResults(exp, seedSteps[1].expectedRes[0])

	// Test intersection end
	exp = s.storage.Events().Select().Intersects(models.Timeframe{
		Start:    time.Date(2022, 1, 12, 13, 30, 0, 0, time.Local),
		Duration: 45 * time.Hour,
	})
	s.ExecuteFilterResults(exp, seedSteps[1].expectedRes[0])
}

func (s *eventsStorageTestSuite) TestFilterUserOperation() {
	s.RunSteps(seedSteps)

	eventIter, err := s.storage.Events().Select().User(2).StartsIn(models.Timeframe{
		Start:    time.Date(2022, 2, 1, 0, 0, 0, 0, time.Local),
		Duration: 30 * 24 * time.Hour,
	}).Execute(context.Background())
	s.NoError(err)

	events, err := eventIter.ToArray()
	s.NoError(err)
	s.Equalf(1, len(events), "Length of filtered events should be %d but was %d", 1, len(events))

	s.Equal(seedSteps[1].expectedRes[1], events[0])
}

func (s *eventsStorageTestSuite) TestFilterMultipleParamsOperation() {
	s.RunSteps(seedSteps)

	exp := s.storage.Events().Select().User(2)
	s.ExecuteFilterResults(exp, seedSteps[1].expectedRes[1])
}

func (s *eventsStorageTestSuite) ExecuteFilterResults(exp basicstorage.EventExpression, e models.Event) {
	eventIter, err := exp.Execute(context.Background())
	s.NoError(err)

	events, err := eventIter.ToArray()
	s.NoError(err)
	s.Equalf(1, len(events), "Length of filtered events should be %d but was %d", 1, len(events))

	s.Equal(e, events[0])
}

func (s *eventsStorageTestSuite) RunSteps(steps []eventStep) {
	for _, step := range steps {
		switch step.action {
		case "create":
			added, err := s.storage.Events().Create(context.Background(), step.e)
			if step.expectedErr != nil {
				s.ErrorIs(err, step.expectedErr)
				continue
			}
			s.Greaterf(added.ID, models.ID(0), "ID of the added element should be more than 0")
			item, e := s.storage.Events().One(context.Background(), added.ID)
			s.NoError(e)
			step.e.ID = added.ID
			s.Equal(step.e, item)
		case "update":
			err := s.storage.Events().Update(context.Background(), step.e)
			if step.expectedErr != nil {
				s.ErrorIs(err, step.expectedErr)
				continue
			}
			item, e := s.storage.Events().One(context.Background(), step.e.ID)
			s.NoError(e)
			s.Equal(step.e, item)
		case "delete":
			err := s.storage.Events().Delete(context.Background(), step.e)
			if step.expectedErr != nil {
				s.ErrorIs(err, step.expectedErr)
				continue
			}
			_, e := s.storage.Events().One(context.Background(), step.e.ID)
			s.Error(e)
		}

		events, e := s.storage.Events().All(context.Background())
		s.NoError(e)

		sort.Slice(events, func(i, j int) bool { return events[i].ID < events[j].ID })

		s.Equal(step.expectedRes, events)
	}
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(eventsStorageTestSuite))
}
