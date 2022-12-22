package client

import (
	"broker-go/internal/broker"
	"broker-go/internal/broker/publisher"
	"broker-go/internal/broker/subscriber"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpClient struct {
	conn *amqp.Connection
}

func NewAmqpClient(conn *amqp.Connection) *AmqpClient {
	return &AmqpClient{conn: conn}
}

func (a *AmqpClient) Close() error {
	return errors.Wrap(a.conn.Close(), "failed to close amqp client")
}

func (a *AmqpClient) Publisher(topic string, opts ...broker.PublishOption) (broker.Publisher, error) {
	ch, err := a.conn.Channel()
	if err != nil {
		return nil, err
	}

	var options broker.PublishOptions
	for _, opt := range opts {
		opt(&options)
	}

	return publisher.NewAmqpPublisher(ch, topic, options.IsPersistent), nil
}

func (a *AmqpClient) Subscriber(topic string, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	ch, err := a.conn.Channel()
	if err != nil {
		return nil, err
	}

	var options broker.SubscribeOptions
	for _, opt := range opts {
		opt(&options)
	}

	if options.IsFanout {
		err := ch.ExchangeDeclare(topic, "fanout", options.IsPersistent, !options.IsPersistent, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	var queueName string = topic
	if options.IsFanout {
		queue, err := ch.QueueDeclare("", options.IsPersistent, !options.IsPersistent, false, false, nil)
		if err != nil {
			return nil, err
		}
		queueName = queue.Name
	} else {
		queue, err := ch.QueueDeclare(topic, options.IsPersistent, !options.IsPersistent, false, false, nil)
		if err != nil {
			return nil, err
		}
		queueName = queue.Name
	}

	if options.IsFanout {
		err := ch.QueueBind(queueName, "", topic, false, nil)
		if err != nil {
			return nil, err
		}
	}

	// log.Printf("%+v", queue)
	return subscriber.NewAmqpSubscriber(ch, queueName), nil
}
