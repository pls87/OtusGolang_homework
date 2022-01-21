package memorystorage

import (
	"context"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
)

type MemoryStorage struct {
	events MemoryEventRepository
	cfg    configs.StorageConf
	mu     *sync.RWMutex
}

func New(cfg configs.StorageConf) *MemoryStorage {
	m := sync.RWMutex{}
	return &MemoryStorage{
		events: MemoryEventRepository{
			mu: &m,
		},
		cfg: cfg,
		mu:  &m,
	}
}

func (s *MemoryStorage) Events() abstractstorage.EventRepository {
	return &s.events
}

func (s *MemoryStorage) Init(_ context.Context) error {
	s.events.Init()

	return nil
}

func (s *MemoryStorage) Dispose() error {
	return nil
}
