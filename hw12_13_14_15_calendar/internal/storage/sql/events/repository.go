package events

// haven't tested this package good enough

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
)

type EventRepository struct {
	DB *sqlx.DB
}

func (s *EventRepository) Init() {
}

func (s *EventRepository) All(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := s.DB.SelectContext(ctx, &events, `SELECT * FROM "events"`)

	return events, err
}

func (s *EventRepository) One(ctx context.Context, id models.ID) (models.Event, error) {
	var ev models.Event
	err := s.DB.GetContext(ctx, &ev, `SELECT * FROM events WHERE ID=$1`, id)

	switch {
	case err == nil:
		return ev, nil
	case errors.Is(err, sql.ErrNoRows):
		return ev, fmt.Errorf("SELECT: event id=%d: %w", id, basic.ErrDoesNotExist)
	default:
		return ev, err
	}
}

func (s *EventRepository) TrackSent(ctx context.Context, eventID models.ID) error {
	query := `INSERT INTO "notification_sent" (event_id, time) VALUES ($1,TIMESTAMP WITH TIME ZONE $2)`
	res, err := s.DB.ExecContext(ctx, query, eventID, time.Now())
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("tracking notification: event id=%d: %w", eventID, basic.ErrNotificationAlreadySent)
	}
	return nil
}

func (s *EventRepository) Create(ctx context.Context, e models.Event) (added models.Event, err error) {
	query := `INSERT INTO "events" (title, user_id, start, duration, notify_before,  description) 
                VALUES ($1, $2, TIMESTAMP WITH TIME ZONE $3, $4, $5, $6) RETURNING "ID"`
	lastID := 0
	err = s.DB.QueryRowxContext(
		ctx, query, e.Title, e.UserID, e.Start, fmt.Sprintf("%d nanoseconds", e.Duration.Nanoseconds()),
		fmt.Sprintf("%d nanoseconds", e.NotifyBefore.Nanoseconds()), e.Desc).Scan(&lastID)

	if err == nil {
		e.ID = models.ID(lastID)
	}

	return e, err
}

func (s *EventRepository) Update(ctx context.Context, e models.Event) error {
	query := `UPDATE "events" SET 
        title=$1, user_id=?, start=TIMESTAMP WITH TIME ZONE $2, 
        duration=$3, notify_before=$4,  description=$5 WHERE ID=$6`
	res, err := s.DB.ExecContext(ctx, query, e.Title, e.UserID, e.Start,
		fmt.Sprintf("%d nanoseconds", e.Duration.Nanoseconds()),
		fmt.Sprintf("%d nanoseconds", e.NotifyBefore.Nanoseconds()), e.Desc, e.ID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("UPDATE: event id=%d: %w", e.ID, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM "events" WHERE ID=$1`, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("DELETE: event id=%d: %w", id, basic.ErrDoesNotExist)
	}
	return nil
}

func (s *EventRepository) DeleteObsolete(ctx context.Context, ttl time.Duration) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM "events" WHERE start + $1 < $2`,
		fmt.Sprintf("%d nanoseconds", ttl.Nanoseconds()), time.Now())
	return err
}

func (s *EventRepository) Select() basic.EventExpression {
	res := EventExpression{
		db:     s.DB,
		params: &basic.EventExpressionParams{},
	}

	return &res
}
