package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/telemetry-platform/alert-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// VehicleEventExchange is the fanout exchange for vehicle lifecycle events.
const VehicleEventExchange = "vehicle_events"

// VehicleEventQueueName is the queue for vehicle deletion events in alert-service.
const VehicleEventQueueName = "alert_vehicle_deletions"

// VehicleEventRetryQueueName is the retry queue for failed vehicle deletion processing.
const VehicleEventRetryQueueName = "alert_vehicle_deletions_retry"

// vehicleDeletedEvent is the payload consumed when a vehicle is deleted.
type vehicleDeletedEvent struct {
	VehicleID int64     `json:"vehicleId"`
	DeletedAt time.Time `json:"deletedAt"`
}

// VehicleDeletionConsumer consumes vehicle.deleted events from RabbitMQ and
// deletes all alerts and the Redis position tracker for the deleted vehicle.
type VehicleDeletionConsumer struct {
	channel  *amqp.Channel
	repo     repositories.AlertRepository
	tracker  resources.PositionTracker
	logger   logger.Logger
}

// VehicleDeletionDelivery is the interface for an AMQP delivery that can be acknowledged.
type VehicleDeletionDelivery interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
}

// Start begins consuming messages from the alert_vehicle_deletions queue.
func (c *VehicleDeletionConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		VehicleEventQueueName, // queue
		"",                    // consumer
		false,                 // autoAck
		false,                 // exclusive
		false,                 // noLocal
		false,                 // noWait
		nil,                   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming vehicle deletion events: %w", err)
	}

	go func() {
		for msg := range msgs {
			c.processMessage(ctx, msg.Body, msg)
		}
	}()

	return nil
}

// Stop closes the consumer channel.
func (c *VehicleDeletionConsumer) Stop() error {
	return c.channel.Close()
}

// processMessage deserializes a vehicle deleted event, deletes alerts and Redis
// tracker, then acks/nacks the message.
func (c *VehicleDeletionConsumer) processMessage(ctx context.Context, body []byte, delivery VehicleDeletionDelivery) {
	var event vehicleDeletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Errorw(ctx, "failed to unmarshal vehicle deleted event", "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	// Delete all alerts for this vehicle from PostgreSQL.
	if err := c.repo.DeleteByVehicleID(ctx, event.VehicleID); err != nil {
		c.logger.Errorw(ctx, "failed to delete alerts for vehicle", "vehicleId", event.VehicleID, "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	// Delete the Redis position tracker key for this vehicle.
	if err := c.tracker.Delete(ctx, event.VehicleID); err != nil {
		c.logger.Errorw(ctx, "failed to delete tracker for vehicle", "vehicleId", event.VehicleID, "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	c.logger.Infow(ctx, "deleted alerts and tracker for vehicle", "vehicleId", event.VehicleID)
	if err := delivery.Ack(false); err != nil {
		c.logger.Errorw(ctx, "failed to ack message", "error", err)
	}
}

// NewVehicleDeletionConsumer creates and returns a new VehicleDeletionConsumer.
// It declares the vehicle_events exchange, the alert_vehicle_deletions queue with
// DLX, the retry queue with TTL and DLX, and binds the queue to the exchange.
func NewVehicleDeletionConsumer(
	client *Client,
	repo repositories.AlertRepository,
	tracker resources.PositionTracker,
	config *env.Configuration,
	log logger.Logger,
) (*VehicleDeletionConsumer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open vehicle deletion consumer channel: %w", err)
	}

	retryTTL, err := time.ParseDuration(config.RetryQueueTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse retry queue TTL: %w", err)
	}

	if err := declareVehicleEventTopology(ch, int(retryTTL.Milliseconds())); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare vehicle event topology: %w", err)
	}

	return &VehicleDeletionConsumer{channel: ch, repo: repo, tracker: tracker, logger: log}, nil
}

// declareVehicleEventTopology declares the vehicle_events fanout exchange, the
// alert_vehicle_deletions queue with DLX args, the retry queue with TTL and DLX args,
// and binds the queue to the exchange.
func declareVehicleEventTopology(ch *amqp.Channel, retryQueueTTLMs int) error {
	if err := ch.ExchangeDeclare(
		VehicleEventExchange, // name
		"fanout",             // kind
		true,                 // durable
		false,                // autoDelete
		false,                // internal
		false,                // noWait
		nil,                  // args
	); err != nil {
		return fmt.Errorf("failed to declare vehicle events exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(
		VehicleEventQueueName, // name
		true,                  // durable
		false,                 // autoDelete
		false,                 // exclusive
		false,                 // noWait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": VehicleEventRetryQueueName,
		},
	); err != nil {
		return fmt.Errorf("failed to declare vehicle deletion queue: %w", err)
	}

	if _, err := ch.QueueDeclare(
		VehicleEventRetryQueueName, // name
		true,                       // durable
		false,                      // autoDelete
		false,                      // exclusive
		false,                      // noWait
		amqp.Table{
			"x-message-ttl":             retryQueueTTLMs,
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": VehicleEventQueueName,
		},
	); err != nil {
		return fmt.Errorf("failed to declare vehicle deletion retry queue: %w", err)
	}

	if err := ch.QueueBind(
		VehicleEventQueueName, // queue
		"",                    // routing key (ignored by fanout)
		VehicleEventExchange,  // exchange
		false,                 // noWait
		nil,                   // args
	); err != nil {
		return fmt.Errorf("failed to bind vehicle deletion queue to exchange: %w", err)
	}

	return nil
}
