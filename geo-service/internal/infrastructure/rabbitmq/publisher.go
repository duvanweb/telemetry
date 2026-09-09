package rabbitmq

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/resources"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Compile-time check that Publisher implements PositionPublisher.
var _ resources.PositionPublisher = (*Publisher)(nil)

// Publisher publishes GPS positions to a RabbitMQ queue.
type Publisher struct {
	channel *amqp.Channel
}

// Publish serializes a position and publishes it to the gps_positions queue.
func (p *Publisher) Publish(ctx context.Context, pos domain.Position) error {
	body, err := json.Marshal(pos)
	if err != nil {
		return fmt.Errorf("failed to marshal position: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		"",           // routing key (ignored by fanout)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish position: %w", err)
	}

	return nil
}

// NewPublisher creates and returns a new Publisher with its own RabbitMQ channel.
func NewPublisher(client *Client, config *env.Configuration) (*Publisher, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open publisher channel: %w", err)
	}

	retryTTL, err := time.ParseDuration(config.RetryQueueTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse retry queue TTL: %w", err)
	}

	if err := DeclareTopology(ch, int(retryTTL.Milliseconds())); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare topology: %w", err)
	}

	return &Publisher{channel: ch}, nil
}
