package logger

import "go.uber.org/fx"

// Module provides the logger dependency via FX.
func Module() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(NewLogger),
	)
}
