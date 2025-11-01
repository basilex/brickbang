package cmd

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"brickbang/internal/config"
	"brickbang/internal/controller"
	"brickbang/internal/middleware"
)

var (
	Version = "none"
	Staging = "none"
	Githash = "none"
	Gobuild = "none"
	Compile = "none"
)

func Metadata() map[string]string {
	return map[string]string{
		"version": Version,
		"staging": Staging,
		"githash": Githash,
		"gobuild": Gobuild,
		"compile": Compile,
	}
}

func ServerPrefork(maxprocs int) int {
	if maxprocs <= 0 {
		maxprocs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxprocs)
	return maxprocs
}

func Run() {
	cfg := config.Get()
	maxprocs := ServerPrefork(cfg.ServerChildProcesses)

	// Server setup
	app := fiber.New(fiber.Config{
		Prefork:               true,
		DisableStartupMessage: true,
		ReadTimeout:           cfg.ServerReadTimeout,
		WriteTimeout:          cfg.ServerWriteTimeout,
		BodyLimit:             cfg.ServerMaxHeaderBytes,
	})

	// CORS middleware
	if cfg.CORSEnabled {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSAllowOrigin,
			AllowMethods:     cfg.CORSAllowMethods,
			AllowHeaders:     cfg.CORSAllowHeaders,
			AllowCredentials: cfg.CORSAllowCredentials,
			ExposeHeaders:    cfg.CORSExposeHeaders,
			MaxAge:           int(cfg.CORSMaxAge.Seconds()),
		}))
	}

	// Middlewares
	app.Use(middleware.AuthMiddleware)
	app.Use(middleware.RBACMiddleware)

	// API versioning
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Route controllers
	controller.RegisterAuxRoutes(v1)
	controller.RegisterUserRoutes(v1)

	// Logging parent and childs
	if !fiber.IsChild() {
		slog.Info("BrickBang server",
			"prefork", true,
			"maxproc", maxprocs,
			"address", cfg.ServerAddress,
			"version", Version,
			"staging", Staging,
			"githash", Githash,
			"gobuild", Gobuild,
			"compile", Compile,
		)
	} else {
		slog.Info(
			"BrickBang child process started", "pid", os.Getpid(), "env", cfg.Env,
		)
	}

	// Server listener
	if err := app.Listen(cfg.ServerAddress); err != nil {
		slog.Error("Error", "Listen", err)
	}
}
