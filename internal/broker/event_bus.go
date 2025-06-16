package broker

import "context"

type EventBusInterface interface {
	Publish(topic Topic, payload []byte, opts ...PublishOption) error
	Subscribe(ctx context.Context, topic Topic, handler func([]byte) error) error
	Close() error
}
