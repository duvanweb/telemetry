package env

import "go.uber.org/fx"

// Module provides the application configuration via FX.
func Module() fx.Option {
	return fx.Module(
		"env",
		fx.Provide(LoadEnv[Configuration]),
	)
}
