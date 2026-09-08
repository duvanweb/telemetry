package services

import (
	"context"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
)

// PositionService is the interface for ingesting GPS positions.
//
//go:generate mockery --name PositionService --dir=. --output=./mocks
type PositionService interface {
	Ingest(ctx context.Context, pos domain.Position) error
}
