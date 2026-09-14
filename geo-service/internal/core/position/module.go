package position

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
)

// Module wires the position domain into FX.
var Module = fx.Module(
	"position",
	fx.Provide(
		fx.Annotate(NewService, fx.As(new(services.PositionService))),
	),
)
