package memorystorage_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/suite"
)

var steps = []basicEventStep{
	{
		action: "create", expectedErr: nil, e: models.Event{
			Timeframe: models.Timeframe{
				Start:    time.Date(2022, 1, 12, 13, 0, 0, 0, time.Local),
				Duration: time.Hour,
			},
			Title:        "Lunch",
			UserID:       1,
			NotifyBefore: 30 * time.Minute,
			Desc:         "Time to eat!",
		}, expectedRes: []models.Event{
			{
				ID: 1,
				Timeframe: models.Timeframe{
					Start:    time.Date(2022, 1, 12, 13, 0, 0, 0, time.Local),
					Duration: time.Hour,
				},
				Title:        "Lunch",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to eat!",
			},
		},
	},
	{
		action: "create", expectedErr: nil, e: models.Event{
			Timeframe: models.Timeframe{
				Start: time.Date(2022, 2, 28, 15, 0, 0, 0,
					time.FixedZone("OMST", 6)),
				Duration: 30 * time.Minute,
			},
			Title:        "Daily Scrum",
			UserID:       1,
			NotifyBefore: 30 * time.Minute,
			Desc:         "Time to meet!",
		}, expectedRes: []models.Event{
			{
				ID: 1,
				Timeframe: models.Timeframe{
					Start:    time.Date(2022, 1, 12, 13, 0, 0, 0, time.Local),
					Duration: time.Hour,
				},
				Title:        "Lunch",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to eat!",
			},
			{
				ID: 2,
				Timeframe: models.Timeframe{
					Start: time.Date(2022, 2, 28, 15, 0, 0, 0,
						time.FixedZone("OMST", 6)),
					Duration: 30 * time.Minute,
				},
				Title:        "Daily Scrum",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet!",
			},
		},
	},
	{
		action: "update", expectedErr: nil, e: models.Event{
			ID: 2,
			Timeframe: models.Timeframe{
				Start: time.Date(2022, 2, 28, 17, 30, 0, 0,
					time.FixedZone("OMST", 6)),
				Duration: 30 * time.Minute,
			},
			Title:        "Daily Scrum",
			UserID:       1,
			NotifyBefore: 30 * time.Minute,
			Desc:         "Time to meet later!",
		}, expectedRes: []models.Event{
			{
				ID: 1,
				Timeframe: models.Timeframe{
					Start:    time.Date(2022, 1, 12, 13, 0, 0, 0, time.Local),
					Duration: time.Hour,
				},
				Title:        "Lunch",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to eat!",
			},
			{
				ID: 2,
				Timeframe: models.Timeframe{
					Start: time.Date(2022, 2, 28, 17, 30, 0, 0,
						time.FixedZone("OMST", 6)),
					Duration: 30 * time.Minute,
				},
				Title:        "Daily Scrum",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet later!",
			},
		},
	},
	{
		action: "update", expectedErr: abstractstorage.ErrDoesNotExist, e: models.Event{
			ID: 3,
		}, expectedRes: []models.Event{
			{
				ID: 1,
				Timeframe: models.Timeframe{
					Start:    time.Date(2022, 1, 12, 13, 0, 0, 0, time.Local),
					Duration: time.Hour,
				},
				Title:        "Lunch",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to eat!",
			},
			{
				ID: 2,
				Timeframe: models.Timeframe{
					Start: time.Date(2022, 2, 28, 17, 30, 0, 0,
						time.FixedZone("OMST", 6)),
					Duration: 30 * time.Minute,
				},
				Title:        "Daily Scrum",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet later!",
			},
		},
	},
	{
		action: "delete", expectedErr: nil, e: models.Event{
			ID: 1,
		}, expectedRes: []models.Event{
			{
				ID: 2,
				Timeframe: models.Timeframe{
					Start: time.Date(2022, 2, 28, 17, 30, 0, 0,
						time.FixedZone("OMST", 6)),
					Duration: 30 * time.Minute,
				},
				Title:        "Daily Scrum",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet later!",
			},
		},
	},
	{
		action: "delete", expectedErr: abstractstorage.ErrDoesNotExist, e: models.Event{
			ID: 1,
		}, expectedRes: []models.Event{
			{
				ID: 2,
				Timeframe: models.Timeframe{
					Start: time.Date(2022, 2, 28, 17, 30, 0, 0,
						time.FixedZone("OMST", 6)),
					Duration: 30 * time.Minute,
				},
				Title:        "Daily Scrum",
				UserID:       1,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet later!",
			},
		},
	},
}

type basicEventStep struct {
	action      string
	e           models.Event
	expectedRes []models.Event
	expectedErr error
}

type eventsStorageTestSuite struct {
	suite.Suite
	storage abstractstorage.Storage
}

func (s *eventsStorageTestSuite) SetupTest() {
	s.storage = memorystorage.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
}

func (s *eventsStorageTestSuite) TearDownTest() {
	s.NoError(s.storage.Dispose())
}

func (s *eventsStorageTestSuite) TestBasicOperations() {
	s.RunSteps(steps)
}

func (s *eventsStorageTestSuite) RunSteps(steps []basicEventStep) {
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

		iter, e := s.storage.Events().All(context.Background())
		s.NoError(e)
		items, e := iter.ToArray()
		s.NoError(e)

		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

		s.Equal(step.expectedRes, items)
	}
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(eventsStorageTestSuite))
}
