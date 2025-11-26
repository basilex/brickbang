package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// Registrar manages the registration of application routes.
type Registrar struct {
	container *Container
	app       *fiber.App
	base      fiber.Router // base group /api/v1
	current   fiber.Router // currently active group
	path      string       // for logging

	// injected middleware instances
	authMiddleware fiber.Handler
	rbacMiddleware fiber.Handler
}

// NewRegistrar initializes the base API structure: /api/v1
func NewRegistrar(
	app *fiber.App,
	container *Container,
	authMiddleware fiber.Handler,
	rbacMiddleware fiber.Handler,
) *Registrar {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	return &Registrar{
		app:            app,
		container:      container,
		path:           "/api/v1",
		base:           v1,
		current:        v1,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

// WithGroup creates a new route group from the base level (/api/v1)
func (rcv *Registrar) WithGroup(prefix string, middlewares ...fiber.Handler) *Registrar {
	group := rcv.base.Group(prefix, middlewares...)
	slog.Debug("Route group registered", "path", rcv.path+prefix)
	rcv.current = group
	return rcv
}

// Public route group (no middleware)
func (rcv *Registrar) WithPublic(prefix string) *Registrar {
	return rcv.WithGroup(prefix)
}

// Private route group (Auth + RBAC middleware)
func (rcv *Registrar) WithPrivate(prefix string) *Registrar {
	return rcv.WithGroup(prefix, rcv.authMiddleware, rcv.rbacMiddleware)
}

// Admin route group
func (rcv *Registrar) WithAdmin(prefix string) *Registrar {
	return rcv.WithGroup(prefix+"/admin", rcv.authMiddleware, rcv.rbacMiddleware)
}

// ========== ROUTE REGISTRATION ==============

func (rcv *Registrar) RegisterAll() *Registrar {

	// PUBLIC
	rcv.WithPublic("/aux").registerAuxRoutes()
	rcv.WithPublic("/auth").registerAuthPublicRoutes()

	// PRIVATE
	rcv.WithPrivate("/auth").registerAuthPrivateRoutes()
	rcv.WithPrivate("/roles").registerRolesRoutes()
	rcv.WithPrivate("/grants").registerGrantsRoutes()
	rcv.WithPrivate("/countries").registerCountryRoutes()
	rcv.WithPrivate("/currencies").registerCurrencyRoutes()

	// ADMIN
	rcv.WithAdmin("/admin").registerSystemAdminRoutes()

	return rcv
}

// =============== MODULES =====================

func (rcv *Registrar) registerAuxRoutes() *Registrar {
	rcv.container.AuxModule.Controller().RegisterRoutes(rcv.current)
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

func (rcv *Registrar) registerRolesRoutes() *Registrar {
	rcv.container.RoleModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerGrantsRoutes() *Registrar {
	rcv.container.GrantModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerCountryRoutes() *Registrar {
	rcv.container.CountryModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerCurrencyRoutes() *Registrar {
	rcv.container.CurrencyModule.Controller().RegisterRoutes(rcv.current)
	return rcv
}

func (rcv *Registrar) registerSystemAdminRoutes() *Registrar {
	// rcv.container.SystemModule.Controller().RegisterAdmin(rcv.current)
	return rcv
}

func (rcv *Registrar) Finalize() {
	// future hooks (metrics, swagger, debug)
}
