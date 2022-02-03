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
	Init() error
	Dispose() error
	Publish(message Message, reliable bool) error
}

type AMPQProducer struct {
	cfg  configs.QueueConf
	conn *amqp.Connection
}

func (ap *AMPQProducer) OpenChannel(reliable bool) (ch *amqp.Channel, err error) {
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

func (ap *AMPQProducer) Init() (err error) {
	ap.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		ap.cfg.User, ap.cfg.Password, ap.cfg.Host, ap.cfg.Port))
	if err != nil {
		return fmt.Errorf("couldn't connect to queue: %w", err)
	}

	var ch *amqp.Channel
	ch, err = ap.conn.Channel()
	if err != nil {
		return fmt.Errorf("couldn't open channel: %w", err)
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(ap.cfg.Exchange, "direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create exchange %s: %w", ap.cfg.Exchange, err)
	}

	if _, err = ch.QueueDeclare(ap.cfg.Queue,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create queue %s: %w", ap.cfg.Queue, err)
	}

	if err = ch.QueueBind(ap.cfg.Queue, ap.cfg.Key, ap.cfg.Exchange, false, nil); err != nil {
		return fmt.Errorf("error binding queue='%s' to exchange='%s' with routing key='%s': %w",
			ap.cfg.Queue, ap.cfg.Exchange, ap.cfg.Key, err)
	}

	if err = ch.Confirm(false); err != nil {
		return fmt.Errorf("error putting channel into confirm mode: %w", err)
	}

	return nil
}

func (ap *AMPQProducer) Dispose() (err error) {
	return ap.conn.Close()
}

func (ap *AMPQProducer) Publish(message Message, reliable bool) (err error) {
	var body []byte
	body, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error while publishing: couldn't marshal message: %w", err)
	}

	var ch *amqp.Channel
	if ch, err = ap.OpenChannel(reliable); err != nil {
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

	return nil
}
