package alert

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/core/ports/services"
)

// Module wires the alert domain into FX.
var Module = fx.Module(
	"alert",
	fx.Provide(
		fx.Annotate(NewService, fx.As(new(services.AlertProcessor))),
	),
)
