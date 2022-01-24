package memory_test

import (
	"time"

	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

var basicSteps = []eventStep{
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
		action: "update", expectedErr: basicstorage.ErrDoesNotExist, e: models.Event{
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
		action: "delete", expectedErr: basicstorage.ErrDoesNotExist, e: models.Event{
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

var seedSteps = []eventStep{
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
			UserID:       2,
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
				UserID:       2,
				NotifyBefore: 30 * time.Minute,
				Desc:         "Time to meet!",
			},
		},
	},
}
