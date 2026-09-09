package services

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// AlertProcessor is the interface for processing positions and detecting anomalies.
//
//go:generate mockery --name AlertProcessor --dir=. --output=./mocks
type AlertProcessor interface {
	Process(ctx context.Context, pos domain.Position) error
}
