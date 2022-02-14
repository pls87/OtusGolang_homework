package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Client
	Consume(tag string) (messages <-chan Message, errors <-chan error, err error)
}

type NotificationConsumer struct {
	NotificationClient
}

func (nc *NotificationConsumer) openChannel() (ch *amqp.Channel, err error) {
	ch, err = nc.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("couldn't open channel: %w", err)
	}

	return ch, err
}

func (nc *NotificationConsumer) Consume(tag string) (messages <-chan Message, errors <-chan error, err error) {
	var ch *amqp.Channel
	if ch, err = nc.openChannel(); err != nil {
		return nil, nil, fmt.Errorf("error while opening channnel for consuming: %w", err)
	}

	var deliveries <-chan amqp.Delivery
	deliveries, err = ch.Consume(
		Queue, // name
		tag,   // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error while consuming messages: %w", err)
	}

	msgs := make(chan Message)
	errs := make(chan error)

	go func() {
		defer func() {
			close(msgs)
			e := ch.Close()
			if e != nil {
				errs <- fmt.Errorf("channel couldn't be closed: %w", e)
			}
			close(errs)
		}()
		var e error
		for d := range deliveries {
			if e = d.Ack(false); e != nil {
				errs <- fmt.Errorf("message couldn't be acknowledged: %w", e)
				continue
			}
			var msg Message
			if e = json.Unmarshal(d.Body, &msg); e != nil {
				errs <- fmt.Errorf("message couldn't be parsed: %w", e)
				continue
			}
			msgs <- msg
		}
	}()

	return msgs, errs, nil
}

func NewConsumer(c configs.QueueConf) Consumer {
	return &NotificationConsumer{
		NotificationClient{
			cfg: c,
		},
	}
}
