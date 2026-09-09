package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/controllers"
)

// Controllers holds all HTTP controllers injected via FX.
type Controllers struct {
	fx.In

	Health *controllers.Health
}

// Router wraps the chi mux and registered routes.
type Router struct {
	controllers Controllers
	server      *chi.Mux
}

// NewRouter creates and returns a new Router.
func NewRouter(server *chi.Mux, c Controllers) *Router {
	return &Router{controllers: c, server: server}
}

// start mounts middlewares and registers all routes under the base path.
func (r *Router) start(basePath string) http.Handler {
	r.server.Use(middleware.RequestID)
	r.server.Use(middleware.Logger)
	r.server.Use(middleware.Recoverer)

	r.server.Get("/health", r.controllers.Health.GetHealth)

	r.server.Route(basePath, func(route chi.Router) {
		// Register business routes here as they are added.
	})

	return r.server
}
