package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/middleware"
)

// Registrator manages the registration of application routes.
// It supports a chain-style grouping and automatic controller registration.
type Registrator struct {
	app       *fiber.App
	container *Container
	base      fiber.Router // base group /api/v1
	current   fiber.Router // currently active group
	path      string       // for logging
}

// NewRegistrator initializes the base API structure: /api/v1
func NewRegistrator(app *fiber.App, container *Container) *Registrator {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	return &Registrator{
		app:       app,
		container: container,
		base:      v1,
		current:   v1,
		path:      "/api/v1",
	}
}

// WithGroup creates a new route group from the base level (/api/v1)
// with optional middleware.
func (r *Registrator) WithGroup(prefix string, middlewares ...fiber.Handler) *Registrator {
	group := r.base.Group(prefix, middlewares...)
	slog.Debug("Route group registered", "path", r.path+prefix)
	r.current = group
	return r
}

// WithPublic defines a public route group (no middleware)
func (r *Registrator) WithPublic(prefix string) *Registrator {
	return r.WithGroup(prefix)
}

// WithPrivate defines a protected route group (Auth + RBAC)
func (r *Registrator) WithPrivate(prefix string) *Registrator {
	return r.WithGroup(prefix, middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// WithAdmin defines an administrative route group
func (r *Registrator) WithAdmin(prefix string) *Registrator {
	return r.WithGroup(prefix+"/admin", middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// RegisterAll registers all application routes
func (r *Registrator) RegisterAll() *Registrator {
	// Public routes
	r.WithPublic("/aux").RegisterAuxRoutes()
	r.WithPublic("/auth").RegisterAuthRoutes()

	// Private routes
	// ...

	return r
}

// RegisterAuxRoutes registers /api/v1/aux routes
func (r *Registrator) RegisterAuxRoutes() *Registrator {
	r.container.AuxController.Register(r.current)
	return r
}

// RegisterAuthRoutes registers /api/v1/auth routes
func (r *Registrator) RegisterAuthRoutes() *Registrator {
	r.container.AuthController.Register(r.current)
	return r
}

// Finalize completes the route registration process (for logging)
func (r *Registrator) Finalize() {
	slog.Info("Routes registration completed", "base_path", r.path)
}
