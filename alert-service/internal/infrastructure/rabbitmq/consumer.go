package rabbitmq

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/services"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Consumer consumes GPS positions from the alert_positions queue and
// delegates processing to the AlertProcessor.
type Consumer struct {
	channel   *amqp.Channel
	processor services.AlertProcessor
	logger    logger.Logger
}

// Delivery is the interface for an AMQP delivery that can be acknowledged.
// amqp.Delivery implements this interface.
type Delivery interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
}

// Start begins consuming messages from the alert_positions queue.
// Each message is deserialized and processed by the AlertProcessor.
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
// It declares the fanout exchange, the alert_positions queue, and binds the queue
// to the exchange.
func NewConsumer(client *Client, processor services.AlertProcessor, log logger.Logger) (*Consumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open consumer channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		ExchangeName, // name
		"fanout",     // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.QueueBind(
		QueueName,    // queue
		"",           // routing key (ignored by fanout)
		ExchangeName, // exchange
		false,        // noWait
		nil,          // args
	); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return &Consumer{channel: ch, processor: processor, logger: log}, nil
}
