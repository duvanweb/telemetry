package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Module wires the HTTP router and server into FX.
func Module() fx.Option {
	return fx.Module(
		"router",
		fx.Provide(
			chi.NewRouter,
			NewRouter,
			controllers.NewHealth,
		),
		fx.Invoke(registerHooks),
	)
}

// registerHooks starts and stops the HTTP server via FX lifecycle.
func registerHooks(lc fx.Lifecycle, shutdown fx.Shutdowner, router *Router, config *env.Configuration, log logger.Logger) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.HTTPPort),
		Handler: router.start("/api"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infow(ctx, "starting http server", "port", config.HTTPPort)
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Errorw(ctx, "http server error", "error", err)
					_ = shutdown.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "stopping http server")
			return server.Shutdown(ctx)
		},
	})
}
