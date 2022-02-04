package notifications

import (
	"fmt"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/streadway/amqp"
)

type Client interface {
	Init() error
	Dispose() error
}

type NotificationClient struct {
	conn *amqp.Connection
	cfg  configs.QueueConf
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

	if err = ch.ExchangeDeclare(nc.cfg.Exchange, "direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create exchange %s: %w", nc.cfg.Exchange, err)
	}

	if _, err = ch.QueueDeclare(nc.cfg.Queue,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create queue %s: %w", nc.cfg.Queue, err)
	}

	if err = ch.QueueBind(nc.cfg.Queue, nc.cfg.Key, nc.cfg.Exchange, false, nil); err != nil {
		return fmt.Errorf("error binding queue='%s' to exchange='%s' with routing key='%s': %w",
			nc.cfg.Queue, nc.cfg.Exchange, nc.cfg.Key, err)
	}

	return nil
}

func (nc *NotificationClient) Dispose() (err error) {
	return nc.conn.Close()
}
