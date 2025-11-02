package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/middleware"
)

// Registrator управляет регистрацией маршрутов приложения.
// Он поддерживает цепочечный стиль группировки и автоматическую регистрацию контроллеров.
type Registrator struct {
	app       *fiber.App
	container *Container
	base      fiber.Router // базовая группа /api/v1
	current   fiber.Router // текущая активная группа
	path      string       // для логирования
}

// NewRegistrator инициализирует базовую структуру API: /api/v1
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

// WithGroup создаёт новую группу маршрутов от базового уровня (/api/v1)
// с возможностью указания middleware.
func (r *Registrator) WithGroup(prefix string, middlewares ...fiber.Handler) *Registrator {
	group := r.base.Group(prefix, middlewares...)
	slog.Debug("Route group registered", "path", r.path+prefix)
	r.current = group
	return r
}

// WithPublic определяет публичную группу маршрутов (без middleware)
func (r *Registrator) WithPublic(prefix string) *Registrator {
	return r.WithGroup(prefix)
}

// WithPrivate определяет защищённую группу маршрутов (Auth + RBAC)
func (r *Registrator) WithPrivate(prefix string) *Registrator {
	return r.WithGroup(prefix, middleware.AuthMiddleware, middleware.RBACMiddleware)
}

// WithAdmin определяет административную группу маршрутов
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

// RegisterAuxRoutes регистрирует маршруты /api/v1/aux
func (r *Registrator) RegisterAuxRoutes() *Registrator {
	r.container.AuxController.Register(r.current)
	return r
}

// RegisterAuthRoutes регистрирует маршруты /api/v1/auth
func (r *Registrator) RegisterAuthRoutes() *Registrator {
	r.container.AuthController.Register(r.current)
	return r
}

// Finalize завершает процесс регистрации маршрутов (для логов)
func (r *Registrator) Finalize() {
	slog.Info("Routes registration completed", "base_path", r.path)
}
