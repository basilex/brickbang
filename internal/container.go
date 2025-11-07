package internal

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"brickbang/internal/config"
	"brickbang/internal/module"

	"brickbang/storage/dbs"
)

// Container holds all dependencies of the application.
// It provides access to shared resources such as database pool and controllers.
type Container struct {
	DBPool    *pgxpool.Pool
	Validator *validator.Validate

	Metadata map[string]string

	// Modules (facades)
	AuxModule     module.IAuxModule
	KeyModule     module.IKeyModule
	AuthModule    module.IAuthModule
	CountryModule module.ICountryModule
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
		panic(err)
	}

	// Initialize metadata
	// and generated SQL queries wrapper (from sqlc)
	metadata := Metadata()
	queries := dbs.New(dbPool)

	// Validator (shared instance)
	validator := validator.New()

	// Initialize modules
	auxModule := module.NewAuxModule(ctx, queries, validator)
	keyModule := module.NewKeyModule(ctx, queries, validator)
	authModule := module.NewAuthModule(ctx, queries, validator)
	countryModule := module.NewCountryModule(ctx, queries, validator)

	// Return a fully initialized dependency container
	return &Container{
		DBPool:        dbPool,
		Metadata:      metadata,
		Validator:     validator,
		AuxModule:     auxModule,
		KeyModule:     keyModule,
		AuthModule:    authModule,
		CountryModule: countryModule,
	}
}
