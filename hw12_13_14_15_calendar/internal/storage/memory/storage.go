package memory

import (
	"context"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
)

type Storage struct {
	events *EventRepository
	cfg    configs.StorageConf
	mu     *sync.RWMutex
}

func New(cfg configs.StorageConf) *Storage {
	m := sync.RWMutex{}
	return &Storage{
		events: &EventRepository{
			mu: &m,
		},
		cfg: cfg,
		mu:  &m,
	}
}

func (s *Storage) Events() basicstorage.EventRepository {
	return s.events
}

func (s *Storage) Init(_ context.Context) error {
	s.events.Init()

	return nil
}

func (s *Storage) Dispose() error {
	return nil
}
