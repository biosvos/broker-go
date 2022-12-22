package message

import amqp "github.com/rabbitmq/amqp091-go"

type AmqpMessage struct {
	delivery *amqp.Delivery
}

func NewAmqpMessage(delivery *amqp.Delivery) *AmqpMessage {
	return &AmqpMessage{delivery: delivery}
}

func (a *AmqpMessage) Body() []byte {
	return a.delivery.Body
}

func (a *AmqpMessage) Ack() error {
	return a.delivery.Ack(false)
}

func (a *AmqpMessage) Nak() error {
	return a.delivery.Nack(false, true)
}
