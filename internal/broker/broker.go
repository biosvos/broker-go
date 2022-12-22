package broker

type Connector interface {
	Connect() (Broker, error)
}

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Close() error

	Publisher(topic string, opts ...PublishOption) (Publisher, error)
	Subscriber(topic string, opts ...SubscribeOption) (Subscriber, error)
}

type Publisher interface {
	Close() error

	Publish(msg []byte) error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Close() error

	Subscribe() (Message, error)
}

type Message interface {
	Body() []byte
	Ack() error
	Nak() error
}
