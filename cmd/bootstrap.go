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
)

// Timeouts are defined
// for startup and shutdown processes timeouts
const (
	startupTimeout  = 500 * time.Millisecond
	shutdownTimeout = 5 * time.Second
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

	// Unified response middleware
	app.Use(middleware.UnifiedResponse())

	// Initialize dependencies and routes
	container := internal.Deps()
	internal.NewRegistrar(app, container).RegisterAll().Finalize()

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
		slog.Info(
			"BrickBang child process started", "pid", os.Getpid(), "env", cfg.Env,
		)
	} else {
		// Master process logs initial info
		slog.Info("BrickBang server starting...")
		slog.Debug("BrickBang server settings:",
			"prefork", true,
			"maxproc", maxprocs,
			"address", cfg.ServerAddress,
			"version", internal.Version,
			"staging", internal.Staging,
			"githash", internal.Githash,
			"gobuild", internal.Gobuild,
			"compile", internal.Compile,
		)

		// Wait a short period to ensure all children started
		time.Sleep(startupTimeout)
		slog.Info("BrickBang server started", "address", cfg.ServerAddress)
	}

	// Wait for shutdown or server error
	select {
	case <-ctx.Done():
		slog.Info("Received shutdown signal, exiting gracefully...")

		// Shutdown Fiber with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("Fiber shutdown error", "error", err)
		}

		// Close DB pool if exists
		if container.DBPool != nil {
			container.DBPool.Close()
			slog.Info("Database pool closed")
		}
	case err := <-serverErr:
		if err != nil {
			slog.Error("Server stopped unexpectedly", "error", err)
		}
	}
	slog.Info("BrickBang server stopped")
}
