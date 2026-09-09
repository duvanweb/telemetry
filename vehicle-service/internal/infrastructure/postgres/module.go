package postgres

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// Module provides the PostgreSQL connection via FX.
func Module() fx.Option {
	return fx.Module(
		"postgres",
		fx.Provide(
			fx.Annotate(NewConnection, fx.As(new(repositories.Databaser))),
		),
		fx.Invoke(registerDBHooks),
	)
}

// registerDBHooks registers the database close hook on the FX lifecycle.
func registerDBHooks(lc fx.Lifecycle, db repositories.Databaser, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "closing database connection")
			return db.Close()
		},
	})
}
