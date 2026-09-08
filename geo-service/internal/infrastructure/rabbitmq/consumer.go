package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Consumer consumes GPS positions from a RabbitMQ queue and persists them.
type Consumer struct {
	channel *amqp.Channel
	repo    repositories.PositionRepository
	logger  logger.Logger
}

// Start begins consuming messages from the gps_positions queue.
// Each message is deserialized and saved via PositionRepository.
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
			c.processMessage(ctx, msg)
		}
	}()

	return nil
}

// Stop closes the consumer channel, stopping message consumption.
func (c *Consumer) Stop() error {
	return c.channel.Close()
}

// processMessage deserializes a position, saves it, and acks/nacks the message.
func (c *Consumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	var pos domain.Position
	if err := json.Unmarshal(msg.Body, &pos); err != nil {
		c.logger.Errorw(ctx, "failed to unmarshal position", "error", err)
		_ = msg.Nack(false, false)
		return
	}

	if err := c.repo.Save(ctx, pos); err != nil {
		c.logger.Errorw(ctx, "failed to save position", "error", err)
		_ = msg.Nack(false, false)
		return
	}

	if err := msg.Ack(false); err != nil {
		c.logger.Errorw(ctx, "failed to ack message", "error", err)
	}
}

// NewConsumer creates and returns a new Consumer with its own RabbitMQ channel.
func NewConsumer(client *Client, repo repositories.PositionRepository, log logger.Logger) (*Consumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open consumer channel: %w", err)
	}

	if _, err := ch.QueueDeclare(QueueName, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{channel: ch, repo: repo, logger: log}, nil
}
