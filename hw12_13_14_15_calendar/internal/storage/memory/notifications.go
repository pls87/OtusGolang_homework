package memory

import (
	"context"
	"sync"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type NotificationRepository struct {
	mu   *sync.RWMutex
	sent map[models.ID]bool
}

func (nr *NotificationRepository) Init() {
	nr.sent = make(map[models.ID]bool)
}

func (nr *NotificationRepository) TrackSent(_ context.Context, eventID models.ID) (err error) {
	nr.mu.Lock()
	defer nr.mu.Unlock()
	if _, ok := nr.sent[eventID]; ok {
		return basic.ErrNotificationAlreadySent
	}
	nr.sent[eventID] = true
	return nil
}
