package broker

type MockRabbitMQPublisher struct {
}

func NewMockRabbitMQPublisher() *MockRabbitMQPublisher {
	return &MockRabbitMQPublisher{}
}

func (m *MockRabbitMQPublisher) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	return nil
}

func (m *MockRabbitMQPublisher) Close() error {
	return nil
}
