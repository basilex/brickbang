package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"brickbang/internal"
	"brickbang/internal/config"
	"brickbang/internal/middleware"
	"brickbang/internal/utility"
)

const (
	startupTimeout = 500 * time.Millisecond
)

// ServerPrefork sets GOMAXPROCS based on the configured child process count.
func ServerPrefork(maxprocs int) int {
	if maxprocs <= 0 || maxprocs > runtime.NumCPU() {
		maxprocs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxprocs)
	return maxprocs
}

// Run initializes the Fiber server, dependencies, routes, and starts listening
func Run() {
	cfg := config.Get()
	maxprocs := ServerPrefork(cfg.ServerChildProcesses)

	// ===== Fiber app setup =====
	app := fiber.New(fiber.Config{
		Prefork:               true,
		DisableStartupMessage: true,
		ReadTimeout:           cfg.ServerReadTimeout,
		WriteTimeout:          cfg.ServerWriteTimeout,
		BodyLimit:             cfg.ServerMaxHeaderBytes,
		ErrorHandler:          middleware.ErrorHandler,
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

	// Default logger middleware
	logger := slog.Default()

	app.Use(middleware.RequestLogger(logger))

	// Global headers response middleware
	app.Use(middleware.GlobalHeaders())

	// Unified response middleware
	app.Use(middleware.UnifiedResponse())

	// ETag siphash middleware
	app.Use(middleware.NewSipHashETagMiddleware(
		cfg.SecuritySipHashKey0, cfg.SecuritySipHashKey1,
	).Handler())

	// Initialize dependencies and routes
	container := internal.NewContainer()

	// Initialize middleware with services
	rbacMW := middleware.RBACMiddleware
	authMW := middleware.NewAuthMiddleware(container.AuthModule.BlacklistService())

	// Initialize Registrar with middleware
	internal.NewRegistrar(app, container, authMW, rbacMW).RegisterAll().Finalize()

	// Signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stop()

	// Channel to capture server errors
	serverErr := make(chan error, 1)

	// Start Fiber server in goroutine
	go func() {
		serverErr <- app.Listen(cfg.ServerAddress)
	}()

	// Master vs Child process logging
	if fiber.IsChild() {
		slog.Info("BrickBang child process started", "pid", os.Getpid(), "env", cfg.Env)
	} else {
		// Master process logs initial info
		slog.Info("BrickBang server starting...")

		if cfg.Env == "dev" && !fiber.IsChild() {
			slog.Info("BrickBang server settings:",
				"prefork", true,
				"maxproc", maxprocs,
				"address", cfg.ServerAddress,
			)

			slog.Info("BrickBang server metadata:",
				"version", internal.Version,
				"staging", internal.Staging,
				"githash", internal.Githash,
				"gobuild", internal.Gobuild,
				"compile", internal.Compile,
			)

			if cfg.Env == "dev" {
				slog.Info("Development mode is enabled")
				utility.InspectRoutes(app, false)
			}
		}

		// TODO: Make it more robust by checking if all children are started
		// Wait a short period to ensure all children started
		time.Sleep(startupTimeout)
		slog.Info("BrickBang server started", "address", cfg.ServerAddress)
	}

	// Wait for shutdown or server error
	select {
	case <-ctx.Done():
		slog.Info("Received shutdown signal, exiting gracefully...")

		// Shutdown Fiber with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("Fiber shutdown error", "error", err)
		}
		closeDatabases(container)
	case err := <-serverErr:
		if err != nil {
			slog.Error("Server stopped unexpectedly", "error", err)
		}
		closeDatabases(container)
	}
	slog.Info("BrickBang server stopped")
}

func closeDatabases(container *internal.Container) {
	// Close Redis Cient if exists
	if container.DBCache != nil {
		container.DBCache.Close()
		slog.Info("Redis connection closed")
	}

	// Close Postgres DB pool if exists
	if container.DBPool != nil {
		container.DBPool.Close()
		slog.Info("Database pool closed")
	}

}
