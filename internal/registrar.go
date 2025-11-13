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
		path:      "/api/v1",
		base:      v1,
		current:   v1,
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
	// ...

	// Public routes
	rcv.WithPublic("/aux").registerAuxRoutes()
	rcv.WithPublic("/auth").registerAuthPublicRoutes()

	// Private routes
	rcv.WithPrivate("/auth").registerAuthPrivateRoutes()
	rcv.WithPrivate("/roles").registerRolesRoutes()

	return rcv
}

// RegisterAuxRoutes registers /api/v1/aux routes
func (rcv *Registrar) registerAuxRoutes() *Registrar {
	rcv.container.AuxModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

// RegisterRoleRoutes registers /api/v1/role routes
func (rcv *Registrar) registerRolesRoutes() *Registrar {
	rcv.container.RoleModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerAuthPublicRoutes() *Registrar {
	rcv.container.AuthModule.Controller().RegisterPublicRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerAuthPrivateRoutes() *Registrar {
	rcv.container.AuthModule.Controller().RegisterPrivateRoutes(rcv.current)
	return rcv
}

// Finalize completes the route registration process (for logging)
func (rcv *Registrar) Finalize() {
	// slog.Info("Routes registration completed", "base_path", rcv.path)
}
