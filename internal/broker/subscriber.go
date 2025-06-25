package broker

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventSubscriber interface {
	Subscribe(ctx context.Context, topic Topic, handler func([]byte) error) error
	Close() error
}

type RabbitMQSubscriber struct {
	url        string
	conn       *amqp.Connection
	subCh      *amqp.Channel
	publisher  EventPublisher
	maxRetries int
}

func NewRabbitMQSubscriber(url string, maxRetries int, publisher EventPublisher) (EventSubscriber, error) {
	s := &RabbitMQSubscriber{
		url:        url,
		publisher:  publisher,
		maxRetries: maxRetries,
	}
	if err := s.connect(); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *RabbitMQSubscriber) connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("subscriber connect: failed to dial: %w", err)
	}
	subCh, err := conn.Channel()
	if err != nil {
		closeConnection(conn)
		return fmt.Errorf("subscriber connect: failed to open channel: %w", err)
	}

	if r.conn != nil {
		closeConnection(r.conn)
	}
	if r.subCh != nil {
		closeChannel(r.subCh)
	}
	r.conn = conn
	r.subCh = subCh
	return nil
}

func (r *RabbitMQSubscriber) Subscribe(ctx context.Context, topic Topic, handler func([]byte) error) error {
	go r.watchConnection(ctx, topic, handler)

	log.Printf("[Rabbit] subscribed to %s", topic)
	return nil
}

func (r *RabbitMQSubscriber) watchConnection(ctx context.Context, topic Topic, handler func([]byte) error) {
	for {
		err := r.consumeMessages(ctx, topic, handler)
		if err != nil {
			log.Printf("Subscribe error (will retry): %v", err)
		}
		select {
		case <-ctx.Done():
			log.Printf("Subscription to %s cancelled", topic)
			return
		default:
			time.Sleep(2 * time.Second)
			if err := r.connect(); err != nil {
				log.Printf("Reconnect failed: %v", err)
				time.Sleep(3 * time.Second)
			}
		}
	}
}

func (r *RabbitMQSubscriber) consumeMessages(ctx context.Context, topic Topic, handler func([]byte) error) error {
	msgs, err := r.subCh.Consume(string(topic), "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}
			r.handleMessage(msg, topic, handler)
		}
	}
}

func (r *RabbitMQSubscriber) handleMessage(msg amqp.Delivery, topic Topic, handler func([]byte) error) {
	retries := r.getRetriesCount(msg.Headers)
	err := handler(msg.Body)
	if err != nil {
		retries++
		if retries > r.maxRetries {
			r.sendToDLQ(msg, topic, retries)
			return
		}
		if errPub := r.publisher.Publish(
			topic,
			msg.Body,
			WithHeaders(amqp.Table{"x-retries": retries}),
			WithContentType(msg.ContentType),
		); errPub != nil {
			log.Printf("Failed to republish message with retry count %d: %v", retries, errPub)
			_ = msg.Nack(false, true)
			return
		}
		_ = msg.Ack(false)
		return
	}
	_ = msg.Ack(false)
}

func (r *RabbitMQSubscriber) getRetriesCount(headers map[string]any) int {
	if hdr, ok := headers[hdrRetries]; ok {
		switch v := hdr.(type) {
		case int32:
			return int(v)
		case int64:
			return int(v)
		case int:
			return v
		case string:
			if val, err := strconv.Atoi(v); err == nil {
				return val
			}
		}
	}
	return 0
}

func (r *RabbitMQSubscriber) sendToDLQ(msg amqp.Delivery, topic Topic, retries int) {
	err := r.publisher.Publish(
		DeadLetterQueue,
		msg.Body,
		WithHeaders(amqp.Table{
			hdrRetries:       retries,
			hdrOriginalTopic: string(topic),
		}),
		WithContentType(msg.ContentType),
	)
	if err != nil {
		log.Printf("DLQ publish failed: %v", err)
		_ = msg.Nack(false, false)
		return
	}
	_ = msg.Ack(false)
	log.Printf("Message sent to DLQ after %d retries", retries)
}

func (r *RabbitMQSubscriber) Close() error {
	var firstErr error
	if r.subCh != nil {
		if err := r.subCh.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
			firstErr = err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}
