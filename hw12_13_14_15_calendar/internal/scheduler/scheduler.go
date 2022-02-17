package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
)

type Scheduler struct {
	prod    notifications.Producer
	storage basic.Storage
	ticker  *time.Ticker
	errors  chan error
}

func NewScheduler(p notifications.Producer, s basic.Storage) Scheduler {
	return Scheduler{
		prod:    p,
		storage: s,
	}
}

func (s *Scheduler) Start(interval time.Duration) (errors <-chan error) {
	s.ticker = time.NewTicker(interval)
	s.errors = make(chan error)
	go func() {
		for range s.ticker.C {
			err := s.removeObsolete(context.Background())
			if err != nil {
				s.errors <- fmt.Errorf("couldn't remove obsolete events: %w", err)
			}
			err = s.generateNotifications()
			if err != nil {
				s.errors <- fmt.Errorf("couldn't generate notifications: %w", err)
			}
		}
	}()
	return s.errors
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	close(s.errors)
}

func (s *Scheduler) removeObsolete(ctx context.Context) error {
	return s.storage.Events().DeleteObsolete(ctx, 24*365*time.Hour)
}

func (s *Scheduler) generateNotifications() error {
	events, err := s.storage.Events().Select().ToNotify().Execute(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get events to notify: %w", err)
	}

	for events.Next() {
		cur, e := events.Current()
		if e != nil {
			return fmt.Errorf("couldn't handle notification: %w", e)
		}
		e = s.prod.Produce(notifications.Message{
			ID:    int(cur.ID),
			Title: cur.Title,
			User:  int(cur.UserID),
			Time:  time.Now(),
		}, false)

		if e != nil {
			return fmt.Errorf("couldn't send notification to queue: %w", e)
		}

		e = s.storage.Notifications().TrackSent(context.Background(), cur.ID)
		if e != nil {
			return fmt.Errorf("couldn't mark notification sent: %w", e)
		}
	}
	return nil
}
