package connector

import (
	"broker-go/internal/broker"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) []byte {
	ret := make([]byte, n)
	for i := range ret {
		ret[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return ret
}

func TestBroker(t *testing.T) {
	conn := NewAmqpConnector("amqp://guest:guest@localhost:5672")
	client, err := conn.Connect()
	assert.NoError(t, err)
	publisher, err := client.Publisher("topic")
	assert.NoError(t, err)

	err = publisher.Publish([]byte("hio"))
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)
}

func TestSubscriber(t *testing.T) {
	conn := NewAmqpConnector("amqp://guest:guest@localhost:5672")
	client, err := conn.Connect()
	assert.NoError(t, err)
	subscriber, err := client.Subscriber("topic")
	assert.NoError(t, err)

	msg, err := subscriber.Subscribe()
	assert.NoError(t, err)
	t.Logf("%s", msg.Body())
	err = msg.Ack()
	assert.NoError(t, err)
}

func TestTwoConsumer(t *testing.T) {
	conn := NewAmqpConnector("amqp://guest:guest@localhost:5672")
	client, err := conn.Connect()
	assert.NoError(t, err)
	defer func() {
		err := client.Close()
		if err != nil {
			t.Log(err)
		}
	}()
	const topic = "two consumer"
	sub1, err := client.Subscriber(topic, broker.WithSubscribePersistent(), broker.WithFanout())
	assert.NoError(t, err)
	defer func() {
		err := sub1.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	sub2, err := client.Subscriber(topic, broker.WithSubscribePersistent(), broker.WithFanout())
	assert.NoError(t, err)
	defer func() {
		err := sub2.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	publisher, err := client.Publisher(topic, broker.WithPublishPersistent())
	assert.NoError(t, err)

	err = publisher.Publish(RandStringRunes(3))
	assert.NoError(t, err)

	msg, err := sub1.Subscribe()
	assert.NoError(t, err)
	log.Println(msg.Body())
	err = msg.Ack()

	msg, err = sub2.Subscribe()
	assert.NoError(t, err)
	log.Println(msg.Body())
	err = msg.Ack()
}
