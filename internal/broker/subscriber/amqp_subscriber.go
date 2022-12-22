package subscriber

import (
	"broker-go/internal/broker"
	"broker-go/internal/broker/message"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpSubscriber struct {
	channel *amqp.Channel
	topic   string
	ch      <-chan amqp.Delivery
}

func NewAmqpSubscriber(channel *amqp.Channel, topic string) *AmqpSubscriber {
	// 2번째 파라미터는 channel 내에서 consumer 들을 identify 하는 필드이다. 만약 "" 라면, 내부에서 자동으로 unique identifier 를 생성한다.
	consume, err := channel.Consume(topic, "", false, false, false, false, nil)
	if err != nil {
		return nil
	}
	return &AmqpSubscriber{
		channel: channel,
		topic:   topic,
		ch:      consume,
	}
}

func (a *AmqpSubscriber) Subscribe() (broker.Message, error) {
	delivery := <-a.ch
	return message.NewAmqpMessage(&delivery), nil
}

func (a *AmqpSubscriber) Close() error {
	return errors.Wrap(a.channel.Close(), "failed to close amqp subscriber")
}
