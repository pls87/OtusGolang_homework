package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
)

type Scheduler struct {
	prod     notifications.Producer
	storage  basic.Storage
	ticker   *time.Ticker
	done     chan interface{}
	eHandler func(e error)
}

func NewScheduler(p notifications.Producer, s basic.Storage, eHandler func(e error)) Scheduler {
	return Scheduler{
		prod:     p,
		storage:  s,
		eHandler: eHandler,
	}
}

func (s *Scheduler) handleTick() {
	e := s.removeObsolete(context.Background())
	if e != nil {
		s.eHandler(fmt.Errorf("couldn't remove obsolete events: %w", e))
	}
	e = s.generateNotifications()
	if e != nil {
		s.eHandler(fmt.Errorf("couldn't generate notifications: %w", e))
	}
}

func (s *Scheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	s.done = make(chan interface{}, 1)
	for {
		select {
		case <-s.ticker.C:
			s.handleTick()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
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

		e = s.storage.Events().TrackSent(context.Background(), cur.ID)
		if e != nil {
			return fmt.Errorf("couldn't mark notification sent: %w", e)
		}
	}
	return nil
}
