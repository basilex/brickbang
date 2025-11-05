package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/middleware"
)

// Registrar manages the registration of application routes.
// It supports a chain-style grouping and automatic controller registration.
type Registrar struct {
	app       *fiber.App
	container *Container
	base      fiber.Router // base group /api/v1
	current   fiber.Router // currently active group
	path      string       // for logging
}

// NewRegistrar initializes the base API structure: /api/v1
func NewRegistrar(app *fiber.App, container *Container) *Registrar {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	return &Registrar{
		app:       app,
		container: container,
		base:      v1,
		current:   v1,
		path:      "/api/v1",
	}
}

// WithGroup creates a new route group from the base level (/api/v1)
// with optional middleware.
func (rcv *Registrar) WithGroup(prefix string, middlewares ...fiber.Handler) *Registrar {
	group := rcv.base.Group(prefix, middlewares...)
	slog.Debug("Route group registered", "path", rcv.path+prefix)
	rcv.current = group
	return rcv
}

// WithPublic defines a public route group (no middleware)
func (rcv *Registrar) WithPublic(prefix string) *Registrar {
	return rcv.WithGroup(prefix)
}

// WithPrivate defines a protected route group (Auth + RBAC)
func (rcv *Registrar) WithPrivate(prefix string) *Registrar {
	return rcv.WithGroup(prefix, middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// WithAdmin defines an administrative route group
func (rcv *Registrar) WithAdmin(prefix string) *Registrar {
	return rcv.WithGroup(prefix+"/admin", middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// RegisterAll registers all application routes
func (rcv *Registrar) RegisterAll() *Registrar {
	// Admin routes
	rcv.WithAdmin("/keys").RegisterKeyRoutes()

	// Public routes
	rcv.WithPublic("/aux").RegisterAuxRoutes()
	rcv.WithPublic("/auth").RegisterAuthRoutes()

	// Private routes
	rcv.WithPrivate("/countries").RegisterCountryRoutes()

	return rcv
}

// RegisterKeyRoutes registers /api/v1/key routes
func (rcv *Registrar) RegisterKeyRoutes() *Registrar {
	rcv.container.KeyController.Register(rcv.current)
	return rcv
}

// RegisterAuxRoutes registers /api/v1/aux routes
func (rcv *Registrar) RegisterAuxRoutes() *Registrar {
	rcv.container.AuxController.Register(rcv.current)
	return rcv
}

// RegisterAuthRoutes registers /api/v1/auth routes
func (rcv *Registrar) RegisterAuthRoutes() *Registrar {
	rcv.container.AuthController.Register(rcv.current)
	return rcv
}

// RegisterCountryRoutes registers /api/v1/countries routes
func (rcv *Registrar) RegisterCountryRoutes() *Registrar {
	rcv.container.CountryController.Register(rcv.current)
	return rcv
}

// Finalize completes the route registration process (for logging)
func (rcv *Registrar) Finalize() {
	slog.Info("Routes registration completed", "base_path", rcv.path)
}
