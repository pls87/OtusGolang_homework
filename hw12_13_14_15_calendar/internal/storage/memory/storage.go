package memorystorage

import (
	"context"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	abstractstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/abstract"
)

type MemoryStorage struct {
	events MemoryEventRepository
	cfg    config.StorageConf
	mu     *sync.RWMutex
}

func New(cfg config.StorageConf) *MemoryStorage {
	m := sync.RWMutex{}
	return &MemoryStorage{
		events: MemoryEventRepository{},
		cfg:    cfg,
		mu:     &m,
	}
}

func (s *MemoryStorage) Events() abstractstorage.EventRepository {
	return &s.events
}

func (s *MemoryStorage) Connect(_ context.Context) error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
