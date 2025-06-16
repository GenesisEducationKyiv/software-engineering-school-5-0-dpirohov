package broker

import (
	"context"
	"log"
	"strconv"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	hdrRetries       = "x-retries"
	hdrOriginalTopic = "x-original-topic"
)

type RabbitMQBus struct {
	conn       *amqp.Connection
	pubCh      *amqp.Channel
	mu         sync.Mutex
	maxRetries int
}

type PublishOption func(*amqp.Publishing)

func NewRabbitMQBus(url string, maxRetries int) (EventBusInterface, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	pubCh, err := conn.Channel()
	if err != nil {
		closeConnection(conn)
		return nil, err
	}

	if err := declareAllQueues(pubCh); err != nil {
		closeChannel(pubCh)
		closeConnection(conn)
		return nil, err
	}

	return &RabbitMQBus{conn: conn, pubCh: pubCh, maxRetries: maxRetries, mu: sync.Mutex{}}, nil
}

func (r *RabbitMQBus) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	pub := amqp.Publishing{
		ContentType:  "application/json",
		Body:         payload,
		DeliveryMode: amqp.Persistent,
	}

	for _, option := range opts {
		option(&pub)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	return r.pubCh.Publish("", string(topic), false, false, pub)
}

func (r *RabbitMQBus) Subscribe(ctx context.Context, topic Topic, handler func([]byte) error) error {
	subCh, err := r.conn.Channel()
	if err != nil {
		log.Printf("open subscribe channel: %v", err)
		return err
	}

	msgs, err := subCh.Consume(string(topic), "", false, false, false, false, nil)
	if err != nil {
		closeChannel(subCh)
		log.Printf("consume: %v", err)
		return err
	}

	go func() {
		defer closeChannel(subCh)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Channel for subscription %s was closed!", topic)
					return
				}
				r.handleMessage(msg, topic, handler)
			}
		}
	}()

	log.Printf("[Rabbit] subscribed to %s", topic)
	return nil
}

func (r *RabbitMQBus) Close() error {
	var firstErr error

	if err := r.pubCh.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ channel: %v", err)
		firstErr = err
	}

	if err := r.conn.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func WithHeaders(h amqp.Table) PublishOption {
	return func(p *amqp.Publishing) { p.Headers = h }
}

func WithContentType(ct string) PublishOption {
	return func(p *amqp.Publishing) { p.ContentType = ct }
}

func (r *RabbitMQBus) getRetriesCount(headers map[string]any) int {
	retries := 0
	if hdr, ok := headers[hdrRetries]; ok {
		if v, ok := hdr.(int32); ok {
			retries = int(v)
		} else if v, ok := hdr.(int64); ok {
			retries = int(v)
		} else if v, ok := hdr.(int); ok {
			retries = v
		} else if v, ok := hdr.(string); ok {
			val, err := strconv.Atoi(v)
			if err == nil {
				retries = val
			} else {
				log.Printf("Failed to convert %s to int\n", v)
			}
		}
	}
	return retries
}

func (r *RabbitMQBus) handleMessage(msg amqp.Delivery, topic Topic, handler func([]byte) error) {
	retries := r.getRetriesCount(msg.Headers)
	err := handler(msg.Body)
	if err != nil {
		retries++
		if retries > r.maxRetries {
			r.sendToDLQ(msg, topic, retries)
			return
		}

		if errPub := r.Publish(
			topic,
			msg.Body,
			WithHeaders(amqp.Table{"x-retries": retries}),
			WithContentType(msg.ContentType),
		); errPub != nil {
			log.Printf("Failed to republish message with retry count %d: %v", retries, errPub)
			if err := msg.Nack(false, true); err != nil {
				log.Printf("Failed to NACK message after publish failure: %v", err)
			}
			return
		}

		if errAck := msg.Ack(false); errAck != nil {
			log.Printf("Failed to ACK old message after republish: %v", errAck)
		}
		return
	}

	if errAck := msg.Ack(false); errAck != nil {
		log.Printf("Failed to ACK message: %v. Body: %s", errAck, string(msg.Body))
	}
}

func (r *RabbitMQBus) sendToDLQ(msg amqp.Delivery, topic Topic, retries int) {
	err := r.Publish(
		DeadLetterQueue,
		msg.Body,
		WithHeaders(amqp.Table{
			hdrRetries:       retries,
			hdrOriginalTopic: string(topic),
		}),
		WithContentType(msg.ContentType),
	)
	if err != nil {
		log.Printf("DLQ publish failed: %v, fallback Nack requeue=false", err)
		if err := msg.Nack(false, false); err != nil {
			log.Printf("Failed to NACK (drop) message: %v. Body: %s", err, string(msg.Body))
		}
		return
	}

	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to ACK after DLQ: %v", err)
	}
	log.Printf("Message sent to DLQ (%s) after %d retries", DeadLetterQueue, retries)
}

func declareAllQueues(ch *amqp.Channel) error {
	for _, t := range AllTopics {
		if _, err := ch.QueueDeclare(string(t), true, false, false, false, nil); err != nil {
			return err
		}
		log.Printf("Queue %s created!", t)
	}
	return nil
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
