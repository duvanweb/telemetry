package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Consumer consumes GPS positions from a RabbitMQ queue and processes them
// through the PositionProcessor (which wraps persistence with a circuit breaker).
type Consumer struct {
	channel   *amqp.Channel
	processor services.PositionProcessor
	logger    logger.Logger
}

// Delivery is the interface for an AMQP delivery that can be acknowledged.
// amqp.Delivery implements this interface.
type Delivery interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
}

// Start begins consuming messages from the gps_positions queue.
// Each message is deserialized and processed via PositionProcessor.
// Runs in a goroutine until the channel is closed.
func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		QueueName, // queue
		"",        // consumer
		false,     // autoAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for msg := range msgs {
			c.processMessage(ctx, msg.Body, msg)
		}
	}()

	return nil
}

// Stop closes the consumer channel, stopping message consumption.
func (c *Consumer) Stop() error {
	return c.channel.Close()
}

// processMessage deserializes a position, processes it, and acks/nacks the message.
// On error (including circuit breaker open), the message is nacked with requeue=false
// so it goes to the retry queue via DLX.
func (c *Consumer) processMessage(ctx context.Context, body []byte, delivery Delivery) {
	var pos domain.Position
	if err := json.Unmarshal(body, &pos); err != nil {
		c.logger.Errorw(ctx, "failed to unmarshal position", "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	if err := c.processor.Process(ctx, pos); err != nil {
		c.logger.Errorw(ctx, "failed to process position", "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	if err := delivery.Ack(false); err != nil {
		c.logger.Errorw(ctx, "failed to ack message", "error", err)
	}
}

// NewConsumer creates and returns a new Consumer with its own RabbitMQ channel.
func NewConsumer(client *Client, processor services.PositionProcessor, config *env.Configuration, log logger.Logger) (*Consumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open consumer channel: %w", err)
	}

	retryTTL, err := time.ParseDuration(config.RetryQueueTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse retry queue TTL: %w", err)
	}

	if err := DeclareTopology(ch, int(retryTTL.Milliseconds())); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare topology: %w", err)
	}

	return &Consumer{channel: ch, processor: processor, logger: log}, nil
}
