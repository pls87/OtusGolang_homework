package events

import (
	"fmt"
	"io"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventIterator struct {
	index int
	items []models.Event
	mu    *sync.RWMutex
}

func (s *EventIterator) Next() bool {
	s.index++
	return s.index < len(s.items)
}

func (s *EventIterator) Current() (models.Event, error) {
	if s.index < len(s.items) {
		return s.items[s.index], nil
	}
	return models.Event{}, fmt.Errorf("iterator is completed: %w", io.EOF)
}

func (s *EventIterator) ToArray() ([]models.Event, error) {
	return s.items, nil
}

func (s *EventIterator) Complete() error {
	return nil
}
