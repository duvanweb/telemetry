package broadcaster

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
)

// Module provides the in-memory AlertBroadcaster via FX.
func Module() fx.Option {
	return fx.Module(
		"broadcaster",
		fx.Provide(
			fx.Annotate(NewBroadcaster, fx.As(new(resources.AlertBroadcaster))),
		),
	)
}
