package health

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/services"
)

// Module wires the health domain into FX.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewService, fx.As(new(services.HealthService))),
	),
)
