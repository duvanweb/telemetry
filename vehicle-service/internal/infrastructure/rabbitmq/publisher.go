package rabbitmq

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/resources"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Compile-time check that Publisher implements resources.EventPublisher.
var _ resources.EventPublisher = (*Publisher)(nil)

// vehicleDeletedEvent is the payload published when a vehicle is deleted.
type vehicleDeletedEvent struct {
	VehicleID int64     `json:"vehicleId"`
	DeletedAt time.Time `json:"deletedAt"`
}

// Publisher publishes vehicle lifecycle events to RabbitMQ.
type Publisher struct {
	channel *amqp.Channel
}

// NewPublisher creates and returns a new Publisher with its own RabbitMQ channel.
// It declares the vehicle_events fanout exchange.
func NewPublisher(client *Client) (*Publisher, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open publisher channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		VehicleEventsExchange, // name
		"fanout",              // kind
		true,                  // durable
		false,                 // autoDelete
		false,                 // internal
		false,                 // noWait
		nil,                   // args
	); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare vehicle events exchange: %w", err)
	}

	return &Publisher{channel: ch}, nil
}

// PublishVehicleDeleted publishes a vehicle.deleted event to the vehicle_events exchange.
func (p *Publisher) PublishVehicleDeleted(ctx context.Context, vehicleID int64, deletedAt time.Time) error {
	event := vehicleDeletedEvent{VehicleID: vehicleID, DeletedAt: deletedAt}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal vehicle deleted event: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		VehicleEventsExchange, // exchange
		"",                    // routing key (ignored by fanout)
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish vehicle deleted event: %w", err)
	}

	return nil
}
