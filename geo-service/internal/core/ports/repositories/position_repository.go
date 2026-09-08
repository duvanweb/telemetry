package repositories

import (
	"context"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
)

// PositionRepository is the port for position persistence.
//
//go:generate mockery --name PositionRepository --dir=. --output=./mocks
type PositionRepository interface {
	Save(ctx context.Context, pos domain.Position) error
}
