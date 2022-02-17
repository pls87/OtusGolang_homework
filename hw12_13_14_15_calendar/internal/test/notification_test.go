package test

import (
	"context"
	"encoding/json"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
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
	return nil
}

func (cp *ConsumerProducer) Consume(tag string) (m <-chan notifications.Message, e <-chan error, err error) {
	var msg notifications.Message
	for s := range cp.N {
		err = json.Unmarshal([]byte(s), &msg)
		if e != nil {
			cp.E <- err
		}
	}
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
}

func (s *NotificationsTestSuite) SetupTest() {
	s.storage = memory.New(configs.StorageConf{})
	s.NoError(s.storage.Init(context.Background()))
	cp := ConsumerProducer{}
	s.NoError(cp.Init())
	s.prod = &cp
	s.con = &cp
}

func (s *NotificationsTestSuite) TearDownTest() {
	s.NoError(s.storage.Dispose())
	s.NoError(s.prod.Dispose())
}
