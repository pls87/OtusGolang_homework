package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func (nr *NotificationRepository) Init() {
}

func (nr *NotificationRepository) TrackSent(ctx context.Context, eventID models.ID) error {
	query := `INSERT INTO "notification_sent" (event_id, time) VALUES ($1,TIMESTAMP WITH TIME ZONE $2)`
	res, err := nr.db.ExecContext(ctx, query, eventID, time.Now())
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("tracking notification: event id=%d: %w", eventID, basic.ErrNotificationAlreadySent)
	}
	return nil
}
