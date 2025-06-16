package broker

import (
	"context"
	"sync"
)

type MockRabbitMQBus struct {
	handlers map[Topic]func([]byte) error
	mu       sync.Mutex
}

func NewMockRabbitMQBus() *MockRabbitMQBus {
	return &MockRabbitMQBus{
		handlers: make(map[Topic]func([]byte) error),
	}
}

func (m *MockRabbitMQBus) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	m.mu.Lock()
	handler, ok := m.handlers[topic]
	m.mu.Unlock()

	if ok {
		return handler(payload)
	}
	return nil
}

func (m *MockRabbitMQBus) Subscribe(ctx context.Context, topic Topic, handler func([]byte) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[topic] = handler
	return nil
}

func (m *MockRabbitMQBus) Close() error {
	return nil
}
