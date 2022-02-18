package sender

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
)

const consumerTag = "calendar_sender"

type Sender struct {
	con      notifications.Consumer
	mHandler func(m notifications.Message)
	eHandler func(e error)
	done     chan interface{}
}

func NewSender(c notifications.Consumer, mhandler func(m notifications.Message), ehandler func(e error)) Sender {
	return Sender{
		con:      c,
		mHandler: mhandler,
		eHandler: ehandler,
	}
}

func (s *Sender) Stop() {
	s.done <- true
}

func (s *Sender) Send() error {
	messages, errors, err := s.con.Consume(consumerTag)
	if err != nil {
		return err
	}
	var e error
	var m notifications.Message
	s.done = make(chan interface{}, 1)
	for ok := true; ok; {
		select {
		case e, ok = <-errors:
			if !ok {
				break
			}
			s.eHandler(e)
		case m, ok = <-messages:
			if !ok {
				break
			}
			s.mHandler(m)
		case <-s.done:
			return nil
		}
	}
	return nil
}
