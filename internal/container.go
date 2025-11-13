package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"brickbang/internal/config"
	"brickbang/internal/module"
	"brickbang/internal/repository/dbs"
)

// Container holds all dependencies of the application.
// It provides access to shared resources such as database pool and controllers.
type Container struct {
	DBPool *pgxpool.Pool

	Metadata map[string]string

	// Modules (facades)
	AuxModule  module.IAuxModule
	AuthModule module.IAuthModule
	RoleModule module.IRoleModule
}

// NewContainer initializes and wires up all application dependencies.
// It follows the dependency injection pattern by creating instances in the correct order:
// Config → Database → Repository → Service → Controller.
func NewContainer() *Container {
	cfg := config.Get()
	ctx := context.Background()

	// Initialize PostgreSQL connection pool
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Container", "failed to initialize database pool", err)
		os.Exit(1)
	}

	// Initialize metadata
	// and generated SQL queries wrapper (from sqlc)
	metadata := Metadata()
	queries := dbs.New(dbPool)

	// Initialize modules
	auxModule := module.NewAuxModule(ctx, metadata)
	authModule := module.NewAuthModule(ctx, queries)
	roleModule := module.NewRoleModule(ctx, queries)

	// Return a fully initialized dependency container
	return &Container{
		DBPool:     dbPool,
		Metadata:   metadata,
		AuxModule:  auxModule,
		AuthModule: authModule,
		RoleModule: roleModule,
	}
}
