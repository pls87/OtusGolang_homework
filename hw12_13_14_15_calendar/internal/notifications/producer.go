package notifications

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

var ErrNotificationWasNotConfirmed = errors.New("publish was not confirmed")

type Producer interface {
	Client
	Publish(message Message, reliable bool) error
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

func (ap *NotificationProducer) Publish(message Message, reliable bool) (err error) {
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

	if err = ch.Publish(ap.cfg.Exchange, ap.cfg.Key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}

	return err
}
