package broker

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	hdrRetries       = "x-retries"
	hdrOriginalTopic = "x-original-topic"
)

type PublishOption func(*amqp.Publishing)

type EventPublisher interface {
	Publish(topic Topic, payload []byte, opts ...PublishOption) error
	Close() error
}

type RabbitMQPublisher struct {
	url         string
	conn        *amqp.Connection
	pubCh       *amqp.Channel
	closeNotify chan *amqp.Error
	mu          sync.Mutex
}

func NewRabbitMQPublisher(url string) (EventPublisher, error) {
	publisher := &RabbitMQPublisher{url: url}
	if err := publisher.connect(); err != nil {
		return nil, err
	}
	go publisher.maintainConnection()
	return publisher, nil
}

func (r *RabbitMQPublisher) connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	pubCh, err := conn.Channel()
	if err != nil {
		closeConnection(conn)
		return fmt.Errorf("channel: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.conn = conn
	r.pubCh = pubCh
	r.closeNotify = conn.NotifyClose(make(chan *amqp.Error, 1))
	return nil
}

func (r *RabbitMQPublisher) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pubCh == nil || r.pubCh.IsClosed() {
		return fmt.Errorf("publish channel is closed")
	}

	pub := amqp.Publishing{
		ContentType:  "application/json",
		Body:         payload,
		DeliveryMode: amqp.Persistent,
	}
	for _, option := range opts {
		option(&pub)
	}
	return r.pubCh.Publish("", string(topic), false, false, pub)
}

func (r *RabbitMQPublisher) maintainConnection() {
	for err := range r.closeNotify {
		log.Printf("RabbitMQ connection closed: %v", err)
		for {
			time.Sleep(time.Second * 5)
			log.Println("Attempting to connect to RabbitMQ...")
			if err := r.connect(); err == nil {
				log.Println("RabbitMQ reconnected successfully.")
				break
			}
			log.Printf("Reconnect failed: %v", err)
		}
	}
}

func (r *RabbitMQPublisher) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var firstErr error

	if r.pubCh != nil {
		if err := r.pubCh.Close(); err != nil {
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

func closeConnection(conn *amqp.Connection) {
	if err := conn.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
	}
}

func closeChannel(ch *amqp.Channel) {
	if err := ch.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ channel: %v", err)
	}
}

func WithHeaders(h amqp.Table) PublishOption {
	return func(p *amqp.Publishing) { p.Headers = h }
}

func WithContentType(ct string) PublishOption {
	return func(p *amqp.Publishing) { p.ContentType = ct }
}
