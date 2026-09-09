package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// QueueName is the name of the RabbitMQ queue for GPS positions.
const QueueName = "gps_positions"

// RetryQueueName is the name of the retry queue for failed position persistence.
const RetryQueueName = "gps_positions_retry"

// ExchangeName is the name of the fanout exchange for GPS positions.
const ExchangeName = "gps_positions_fanout"

// Client holds the RabbitMQ connection.
type Client struct {
	connection *amqp.Connection
}

// NewClient creates and returns a new RabbitMQ client.
// It dials the RabbitMQ server and verifies connectivity.
func NewClient(config *env.Configuration, log logger.Logger) (*Client, error) {
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	log.Infow(context.Background(), "rabbitmq connection established", "url", config.RabbitMQURL)

	return &Client{connection: conn}, nil
}

// Channel opens a new channel on the RabbitMQ connection.
func (c *Client) Channel() (*amqp.Channel, error) {
	return c.connection.Channel()
}

// Close closes the RabbitMQ connection.
func (c *Client) Close() error {
	return c.connection.Close()
}

// DeclareTopology declares the fanout exchange, the gps_positions queue with DLX
// args, the gps_positions_retry queue with TTL and DLX args, and binds the
// gps_positions queue to the exchange. retryQueueTTLMs is the TTL for the retry
// queue in milliseconds.
func DeclareTopology(ch *amqp.Channel, retryQueueTTLMs int) error {
	if err := ch.ExchangeDeclare(
		ExchangeName, // name
		"fanout",     // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": RetryQueueName,
		},
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if _, err := ch.QueueDeclare(
		RetryQueueName, // name
		true,           // durable
		false,          // autoDelete
		false,          // exclusive
		false,          // noWait
		amqp.Table{
			"x-message-ttl":             retryQueueTTLMs,
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": QueueName,
		},
	); err != nil {
		return fmt.Errorf("failed to declare retry queue: %w", err)
	}

	if err := ch.QueueBind(
		QueueName,    // queue
		"",           // routing key (ignored by fanout)
		ExchangeName, // exchange
		false,        // noWait
		nil,          // args
	); err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return nil
}
