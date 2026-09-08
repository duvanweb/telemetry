package resources

import (
	"context"
	"time"
)

// PositionCache is the port for the anti-duplicate cache.
//
//go:generate mockery --name PositionCache --dir=. --output=./mocks
type PositionCache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}
