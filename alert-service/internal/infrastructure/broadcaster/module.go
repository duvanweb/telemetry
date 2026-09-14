package broadcaster

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
)

// Module provides the in-memory AlertBroadcaster and PositionBroadcaster via FX.
func Module() fx.Option {
	return fx.Module(
		"broadcaster",
		fx.Provide(
			fx.Annotate(NewBroadcaster, fx.As(new(resources.AlertBroadcaster))),
			fx.Annotate(NewPositionBroadcaster, fx.As(new(resources.PositionBroadcaster))),
		),
	)
}
