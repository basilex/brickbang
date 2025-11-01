package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/middleware"
)

// Registrator manages route registration with chainable groups.
type Registrator struct {
	app       *fiber.App
	container *Container
	current   fiber.Router
	path      string
}

// NewRegistrator initializes base API groups (/api/v1)
func NewRegistrator(app *fiber.App, container *Container) *Registrator {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	return &Registrator{
		app:       app,
		container: container,
		current:   v1,
		path:      "/api/v1",
	}
}

// WithGroup creates a new nested route group with optional middlewares
func (r *Registrator) WithGroup(prefix string, middlewares ...fiber.Handler) *Registrator {
	group := r.current.Group(prefix)

	// Convert []fiber.Handler to []interface{} for group.Use()
	if len(middlewares) > 0 {
		handlers := make([]interface{}, len(middlewares))
		for i, h := range middlewares {
			handlers[i] = h
		}
		group.Use(handlers...)
	}
	r.current = group
	r.path += prefix
	slog.Debug("Route group registered", "path", r.path)
	return r
}

// WithPublic defines a public route group
func (r *Registrator) WithPublic() *Registrator {
	return r.WithGroup("") // /api/v1
}

// WithPrivate defines a protected route group (Auth + RBAC)
func (r *Registrator) WithPrivate() *Registrator {
	return r.WithGroup("", middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// WithAdmin defines an admin route group
func (r *Registrator) WithAdmin() *Registrator {
	return r.WithGroup("/admin", middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// RegisterAll registers all controllers from the container
func (r *Registrator) RegisterAll() *Registrator {
	r.WithPublic().WithGroup("/aux").RegisterAuxRoutes()
	r.WithPrivate().WithGroup("/user").RegisterUserRoutes()
	return r
}

// RegisterAuxRoutes registers /aux endpoints
func (r *Registrator) RegisterAuxRoutes() *Registrator {
	r.container.AuxController.Register(r.current)
	return r
}

// RegisterUserRoutes registers /user endpoints
func (r *Registrator) RegisterUserRoutes() *Registrator {
	r.container.AuthController.Register(r.current)
	return r
}

// Finalize completes registration and can log final path
func (r *Registrator) Finalize() {
	// slog.Info("Routes registration completed", "base_path", r.path)
}
