package connector

import (
	"broker-go/internal/broker"
	"broker-go/internal/broker/client"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpConnector struct {
	url string
}

func NewAmqpConnector(url string) *AmqpConnector {
	return &AmqpConnector{url: url}
}

func (c *AmqpConnector) Connect() (broker.Broker, error) {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return nil, err
	}
	return client.NewAmqpClient(conn), nil
}
