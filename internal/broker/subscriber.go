package broker

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventSubscriber interface {
	Subscribe(ctx context.Context, topic Topic, handler func(ctx context.Context, data []byte) error) error
	Close() error
}

type RabbitMQSubscriber struct {
	url        string
	conn       *amqp.Connection
	subCh      *amqp.Channel
	publisher  EventPublisher
	maxRetries int
}

func NewRabbitMQSubscriber(url string, maxRetries int, publisher EventPublisher) (*RabbitMQSubscriber, error) {
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

func (r *RabbitMQSubscriber) Subscribe(ctx context.Context, topic Topic, handler func(ctx context.Context, data []byte) error) error {
	if err := r.connect(); err != nil {
		logger.Log.Fatal().Err(err).Msgf("RabbitMQSubscriber: failed to connectd to %s!", topic)
	}

	if err := r.declareQueues(topic); err != nil {
		closeChannel(r.subCh)
		closeConnection(r.conn)
		return fmt.Errorf("error while declaring queues: %w", err)
	}

	go r.maintainConnection(ctx, topic, handler)

	logger.Log.Info().Msgf("[Rabbit] subscribed to %s", topic)
	return nil
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

func (r *RabbitMQSubscriber) maintainConnection(ctx context.Context, topic Topic, handler func(ctx context.Context, data []byte) error) {
	for {
		err := r.consumeMessages(ctx, topic, handler)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Subscribe error, retrying...")
		}
		select {
		case <-ctx.Done():
			logger.Log.Info().Msgf("Subscription to %s cancelled", topic)
			return
		default:
			time.Sleep(2 * time.Second)
			if err := r.connect(); err != nil {
				logger.Log.Error().Err(err).Msg("Reconnect failed")
				time.Sleep(3 * time.Second)
			}
		}
	}
}

func (r *RabbitMQSubscriber) consumeMessages(ctx context.Context, topic Topic, handler func(ctx context.Context, data []byte) error) error {
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
			r.handleMessage(ctx, msg, topic, handler)
		}
	}
}

func (r *RabbitMQSubscriber) handleMessage(ctx context.Context, msg amqp.Delivery, topic Topic, handler func(ctx context.Context, data []byte) error) {
	retries := r.getRetriesCount(msg.Headers)
	ctx = logger.InjectTraceID(ctx, msg)
	log := logger.FromContext(ctx)
	err := handler(ctx, msg.Body)
	if err != nil {
		retries++
		if retries > r.maxRetries {
			r.sendToDLQ(ctx, msg, topic, retries)
			return
		}
		if errPub := r.publisher.Publish(
			topic,
			msg.Body,
			WithHeaders(amqp.Table{"x-retries": retries}),
			WithContentType(msg.ContentType),
		); errPub != nil {
			log.Error().Err(errPub).Msgf("Failed to republish message with retry count %d", retries)
			_ = msg.Nack(false, true)
			return
		}
		_ = msg.Ack(false)
		return
	}
	_ = msg.Ack(false)
}

func (r *RabbitMQSubscriber) getRetriesCount(headers map[string]any) int {
	if hdr, ok := headers[constants.HdrRetries]; ok {
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

func (r *RabbitMQSubscriber) sendToDLQ(ctx context.Context, msg amqp.Delivery, topic Topic, retries int) {
	log := logger.FromContext(ctx)

	err := r.publisher.Publish(
		topic.DLQ(),
		msg.Body,
		WithHeaders(amqp.Table{
			constants.HdrRetries:       retries,
			constants.HdrOriginalTopic: string(topic),
		}),
		WithContentType(msg.ContentType),
	)
	if err != nil {
		log.Error().Err(err).Msg("DLQ publish failed")
		_ = msg.Nack(false, false)
		return
	}
	_ = msg.Ack(false)
	log.Info().Msgf("Message sent to DLQ after %d retries", retries)
}

func (r *RabbitMQSubscriber) Close() error {
	var firstErr error
	if r.subCh != nil {
		if err := r.subCh.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ channel")
			firstErr = err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ connection")
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (r *RabbitMQSubscriber) declareQueues(topic Topic) error {
	for _, t := range []Topic{topic, topic.DLQ()} {
		if _, err := r.subCh.QueueDeclare(string(t), true, false, false, false, nil); err != nil {
			return err
		}
		logger.Log.Info().Msg("Queue %s created!")
	}
	return nil
}
