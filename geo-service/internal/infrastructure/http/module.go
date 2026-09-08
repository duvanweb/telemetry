package httpclient

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/ports/resources"
)

// Module provides the vehicle-service HTTP client via FX.
func Module() fx.Option {
	return fx.Module(
		"httpclient",
		fx.Provide(
			fx.Annotate(NewVehicleClient, fx.As(new(resources.VehicleClient))),
		),
	)
}
