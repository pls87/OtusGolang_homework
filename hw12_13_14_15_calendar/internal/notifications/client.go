package notifications

import (
	"fmt"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/streadway/amqp"
)

const (
	Exchange = "calendar"
	Queue    = "notifications"
	Key      = "new_notification"
)

type Client interface {
	Init() error
	Dispose() error
}

type NotificationClient struct {
	conn *amqp.Connection
	cfg  configs.NotificationConf
}

func (nc *NotificationClient) Init() (err error) {
	nc.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		nc.cfg.User, nc.cfg.Password, nc.cfg.Host, nc.cfg.Port))
	if err != nil {
		return fmt.Errorf("couldn't connect to queue: %w", err)
	}

	var ch *amqp.Channel
	ch, err = nc.conn.Channel()
	if err != nil {
		return fmt.Errorf("couldn't open channel: %w", err)
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(Exchange, "direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create exchange %s: %w", Exchange, err)
	}

	if _, err = ch.QueueDeclare(Queue,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create queue %s: %w", Queue, err)
	}

	if err = ch.QueueBind(Queue, Key, Exchange, false, nil); err != nil {
		return fmt.Errorf("error binding queue='%s' to exchange='%s' with routing key='%s': %w",
			Queue, Exchange, Key, err)
	}

	return nil
}

func (nc *NotificationConsumer) openChannel() (ch *amqp.Channel, err error) {
	if nc.conn == nil {
		return nil, fmt.Errorf("connection is not opened: %w", ErrCouldNotOpenChannel)
	}
	ch, err = nc.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("couldn't open channel: %w", err)
	}

	return ch, err
}

func (nc *NotificationClient) Dispose() (err error) {
	if nc.conn != nil {
		return nc.conn.Close()
	}
	return nil
}
