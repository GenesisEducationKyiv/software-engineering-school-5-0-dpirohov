package broker

import (
	"context"
)

type MockEventSubscriber struct {
	SubscribedTopics []Topic
	Handlers         map[Topic]func(ctx context.Context, data []byte) error
}

func NewMockEventSubscriber() *MockEventSubscriber {
	return &MockEventSubscriber{
		Handlers: make(map[Topic]func(ctx context.Context, data []byte) error),
	}
}

func (m *MockEventSubscriber) Subscribe(ctx context.Context, topic Topic, handler func(ctx context.Context, data []byte) error) error {
	m.SubscribedTopics = append(m.SubscribedTopics, topic)
	m.Handlers[topic] = handler
	return nil
}

func (m *MockEventSubscriber) Close() error {
	return nil
}

func (m *MockEventSubscriber) SimulateMessage(ctx context.Context, topic Topic, payload []byte) error {
	if handler, ok := m.Handlers[topic]; ok {
		return handler(ctx, payload)
	}
	return nil
}
