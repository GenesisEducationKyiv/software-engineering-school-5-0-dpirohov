package broker

type MockRabbitMQPublisher struct {
	Calls []PublishCall
}

type PublishCall struct {
	Topic   Topic
	Payload []byte
	Options []PublishOption
}

func NewMockRabbitMQPublisher() *MockRabbitMQPublisher {
	return &MockRabbitMQPublisher{}
}

func (m *MockRabbitMQPublisher) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	m.Calls = append(m.Calls, PublishCall{
		Topic:   topic,
		Payload: payload,
		Options: opts,
	})
	return nil
}

func (m *MockRabbitMQPublisher) Close() error {
	return nil
}
