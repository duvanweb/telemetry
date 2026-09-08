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
