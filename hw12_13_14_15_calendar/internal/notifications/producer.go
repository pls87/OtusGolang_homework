package notifications

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/streadway/amqp"
)

var ErrNotificationWasNotConfirmed = errors.New("publish was not confirmed")

type Producer interface {
	Client
	Produce(message Message, reliable bool) error
}

type NotificationProducer struct {
	NotificationClient
}

func (ap *NotificationProducer) openChannel(reliable bool) (ch *amqp.Channel, err error) {
	ch, err = ap.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("couldn't open channel: %w", err)
	}

	if reliable {
		if err = ch.Confirm(false); err != nil {
			return nil, fmt.Errorf("error while putting channel into confirm mode: %w", err)
		}
	}

	return ch, err
}

func (ap *NotificationProducer) Produce(message Message, reliable bool) (err error) {
	var body []byte
	body, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error while publishing: couldn't marshal message: %w", err)
	}

	var ch *amqp.Channel
	if ch, err = ap.openChannel(reliable); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}
	defer ch.Close()

	if reliable {
		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer func() {
			if confirmed := <-confirms; !confirmed.Ack && err == nil {
				err = fmt.Errorf("error while publishing: %w", ErrNotificationWasNotConfirmed)
			}
		}()
	}

	if err = ch.Publish(Exchange, Key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}

	return err
}

func NewProducer(c configs.NotificationConf) Producer {
	return &NotificationProducer{
		NotificationClient{cfg: c},
	}
}
