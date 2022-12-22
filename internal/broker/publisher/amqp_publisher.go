package publisher

import (
	"context"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpPublisher struct {
	ch           *amqp.Channel
	topic        string
	deliveryMode uint8
}

func (a *AmqpPublisher) Close() error {
	return errors.Wrap(a.ch.Close(), "failed to close amqp publisher")
}

func NewAmqpPublisher(ch *amqp.Channel, topic string, isPersist bool) *AmqpPublisher {
	ret := AmqpPublisher{
		ch:    ch,
		topic: topic,
	}
	if isPersist {
		ret.deliveryMode = amqp.Persistent
	} else {
		ret.deliveryMode = amqp.Transient
	}
	return &ret
}

func (a *AmqpPublisher) Publish(msg []byte) error {
	return a.ch.PublishWithContext(context.Background(), a.topic, "", false, false, amqp.Publishing{
		DeliveryMode: a.deliveryMode,
		Body:         msg,
	})
}
