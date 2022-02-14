package basic

import (
	"context"
	"errors"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

var ErrNotificationAlreadySent = errors.New("notification already sent")

type NotificationRepository interface {
	Init()
	TrackSent(ctx context.Context, ID models.ID) error
}
