package broker

import (
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

type RabbitMQBus struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type PublishOption func(*amqp.Publishing)

func NewRabbitMQBus(url string) (EventBusInerface, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQBus{conn: conn, ch: ch}, nil
}

func (r *RabbitMQBus) Publish(topic Topic, payload []byte, opts ...PublishOption) error {
	_, err := r.ch.QueueDeclare(string(topic), true, false, false, false, nil)
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		ContentType:  "application/json",
		Body:         payload,
		DeliveryMode: amqp.Persistent,
	}

	for _, option := range opts {
		option(&pub)
	}

	return r.ch.Publish("", string(topic), false, false, pub)
}

func (r *RabbitMQBus) Subscribe(topic Topic, handler func([]byte) error) error {
	_, err := r.ch.QueueDeclare(string(topic), true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := r.ch.Consume(string(topic), "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	const maxRetries = 3

	go func() {
		for msg := range msgs {
			r.handleMessage(msg, topic, handler, maxRetries)
		}
	}()

	return nil
}

func (r *RabbitMQBus) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

func WithHeaders(h amqp.Table) PublishOption {
	return func(p *amqp.Publishing) { p.Headers = h }
}

func WithContentType(ct string) PublishOption {
	return func(p *amqp.Publishing) { p.ContentType = ct }
}

func (r *RabbitMQBus) getRetriesCount(headers map[string]any) int {
	retries := 0
	if hdr, ok := headers["x-retries"]; ok {
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
				log.Println("Failed to convert &s to int", v)
			}
		}
	}
	return retries
}

func (r *RabbitMQBus) handleMessage(msg amqp.Delivery, topic Topic, handler func([]byte) error, maxRetries int) {
	retries := r.getRetriesCount(msg.Headers)
	err := handler(msg.Body)
	if err != nil {
		retries++
		if retries > maxRetries {
			if err := msg.Nack(false, false); err != nil {
				log.Printf("Failed to NACK (drop) message: %v. Body: %s", err, string(msg.Body))
			} else {
				log.Printf("Message dropped after %d retries. Body: %s", retries, string(msg.Body))
			}
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
