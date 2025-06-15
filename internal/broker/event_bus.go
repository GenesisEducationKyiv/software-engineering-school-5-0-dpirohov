package broker

type EventBusInerface interface {
	Publish(topic Topic, payload []byte, opts ...PublishOption) error
	Subscribe(topic Topic, handler func([]byte) error) error
	Close() error
}
