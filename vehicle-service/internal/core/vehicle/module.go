package vehicle

import "go.uber.org/fx"

// Module wires the vehicle domain into FX.
var Module = fx.Module(
	"vehicle",
	fx.Provide(NewService),
)
