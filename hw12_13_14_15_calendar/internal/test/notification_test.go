package test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/sender"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/suite"
)

type ConsumerProducer struct { // котопес
	N chan string
	E chan error
	M chan notifications.Message
}

func (cp *ConsumerProducer) Init() (err error) {
	cp.N = make(chan string)
	cp.E = make(chan error)
	cp.M = make(chan notifications.Message)
	return nil
}

func (cp *ConsumerProducer) Dispose() (err error) {
	close(cp.N)
	close(cp.E)
	close(cp.M)
	return nil
}

func (cp *ConsumerProducer) Consume(_ string) (m <-chan notifications.Message, errors <-chan error, err error) {
	go func() {
		for s := range cp.N {
			var msg notifications.Message
			e := json.Unmarshal([]byte(s), &msg)
			if e != nil {
				cp.E <- e
			} else {
				cp.M <- msg
			}
		}
	}()
	return cp.M, cp.E, nil
}

func (cp *ConsumerProducer) Produce(message notifications.Message, reliable bool) (err error) {
	bytes, e := json.Marshal(message)
	if e != nil {
		return e
	}
	cp.N <- string(bytes)
	return nil
}

type NotificationsTestSuite struct {
	suite.Suite
	storage basic.Storage
	prod    notifications.Producer
	con     notifications.Consumer

	scheduler scheduler.Scheduler
	sender    sender.Sender

	mu            *sync.Mutex
	notifications []notifications.Message
}

func (s *NotificationsTestSuite) SetupSuite() {
	s.mu = &sync.Mutex{}
	cp := ConsumerProducer{}
	s.NoError(cp.Init())
	s.prod = &cp
	s.con = &cp
	s.storage = memory.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
	s.scheduler = scheduler.NewScheduler(s.prod, s.storage, s.onError)
	s.sender = sender.NewSender(s.con, s.onMessage, s.onError)
	go func() {
		s.scheduler.Start(100 * time.Millisecond)
	}()
	go func() {
		s.NoError(s.sender.Send())
	}()
}

func (s *NotificationsTestSuite) SetupTest() {
	s.notifications = make([]notifications.Message, 0, 10)
}

func (s *NotificationsTestSuite) TearDownSuite() {
	s.scheduler.Stop()
	s.sender.Stop()
	s.NoError(s.prod.Dispose())
	s.NoError(s.storage.Dispose())
}

func (s *NotificationsTestSuite) findMessage(id int) []notifications.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]notifications.Message, 0, 3)
	for _, v := range s.notifications {
		if v.ID == id {
			res = append(res, v)
		}
	}
	return res
}

func (s *NotificationsTestSuite) onMessage(m notifications.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications = append(s.notifications, m)
}

func (s *NotificationsTestSuite) onError(e error) {
	s.Failf("no error should be during this test but happened: %s", e.Error())
}

func (s *NotificationsTestSuite) TestNotificationShouldCome() {
	added, err := s.storage.Events().Create(context.Background(), models.Event{
		Timeframe: models.Timeframe{
			Start:    time.Now().Add(5 * time.Minute),
			Duration: 30 * time.Minute,
		},
		Title:        "Morning Coffee",
		UserID:       1,
		NotifyBefore: 15 * time.Minute,
	})

	s.NoError(err)

	s.Eventuallyf(func() bool {
		n := s.findMessage(int(added.ID))
		return len(n) > 0
	}, 150*time.Millisecond, 10*time.Millisecond, "notification should come")

	s.Eventuallyf(func() bool {
		n := s.findMessage(int(added.ID))
		return len(n) == 1
	}, 200*time.Millisecond, 150*time.Millisecond, "notification should come once")
}

func (s *NotificationsTestSuite) TestNotificationShouldNotCome() {
	added, err := s.storage.Events().Create(context.Background(), models.Event{
		Timeframe: models.Timeframe{
			Start:    time.Now().Add(1 * time.Hour),
			Duration: 30 * time.Minute,
		},
		Title:        "Morning Coffee",
		UserID:       1,
		NotifyBefore: 59 * time.Minute,
	})

	s.NoError(err)

	s.Eventuallyf(func() bool {
		n := s.findMessage(int(added.ID))
		return len(n) == 0
	}, 200*time.Millisecond, 150*time.Millisecond, "notification should not come")
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(NotificationsTestSuite))
}
