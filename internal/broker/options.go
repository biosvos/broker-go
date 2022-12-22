package broker

type PublishOptions struct {
	IsPersistent bool
}

type PublishOption func(*PublishOptions)

func WithPublishPersistent() PublishOption {
	return func(o *PublishOptions) {
		o.IsPersistent = true
	}
}

type SubscribeOptions struct {
	IsPersistent bool
	IsFanout     bool
}

type SubscribeOption func(*SubscribeOptions)

func WithSubscribePersistent() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.IsPersistent = true
	}
}

func WithFanout() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.IsFanout = true
	}
}
