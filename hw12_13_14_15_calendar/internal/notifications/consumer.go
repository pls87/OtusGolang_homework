package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Client
	Consume(tag string) (messages chan Message, errors chan error, err error)
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

func (nc *NotificationConsumer) Consume(tag string) (messages chan Message, errors chan error, err error) {
	var ch *amqp.Channel
	if ch, err = nc.openChannel(); err != nil {
		return nil, nil, fmt.Errorf("error while publishing: %w", err)
	}

	deliveries, err := ch.Consume(
		nc.cfg.Queue, // name
		tag,          // consumerTag,
		false,        // noAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error while consuming messages: %s", err)
	}

	messages = make(chan Message)
	errors = make(chan error)

	go func() {
		defer func() {
			close(messages)
			e := ch.Close()
			if e != nil {
				errors <- fmt.Errorf("channel couldn't be closed: %w", e)
			}
			close(errors)
		}()
		var e error
		for d := range deliveries {
			if e = d.Ack(false); e != nil {
				errors <- fmt.Errorf("message couldn't be acknowledged: %w", e)
				continue
			}
			var msg Message
			if e = json.Unmarshal(d.Body, &msg); e != nil {
				errors <- fmt.Errorf("message couldn't be parsed: %w", e)
				continue
			}
			messages <- msg
		}
	}()

	return messages, errors, nil
}
