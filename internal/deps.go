package internal

import (
	"context"

	"brickbang/internal/config"
	"brickbang/internal/controller"
	"brickbang/internal/repository"
	"brickbang/internal/service"
	"brickbang/storage/dbs"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Container holds all dependencies of the application.
// It provides access to shared resources such as database pool and controllers.
type Container struct {
	DBPool *pgxpool.Pool

	AuxController  controller.IAuxController
	AuthController controller.IAuthController
}

// Deps initializes and wires up all application dependencies.
// It follows the dependency injection pattern by creating instances in the correct order:
// Config → Database → Repository → Service → Controller.
func Deps() *Container {
	cfg := config.Get()
	ctx := context.Background()

	// Initialize PostgreSQL connection pool
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	// Initialize generated SQL queries wrapper (from sqlc)
	queries := dbs.New(dbPool)

	// Repository layer: responsible for database access
	auxRepo := repository.NewAuxRepository()
	authRepo := repository.NewAuthRepository(queries)

	// Service layer: contains business logic
	auxService := service.NewAuxService(auxRepo)
	authService := service.NewAuthService(authRepo, cfg.JWTSecret, cfg.JWTAccessExpiration)

	// Controller layer: handles HTTP requests/responses
	auxController := controller.NewAuxController(auxService)
	authController := controller.NewAuthController(authService)

	// Return a fully initialized dependency container
	return &Container{
		DBPool:         dbPool,
		AuxController:  auxController,
		AuthController: authController,
	}
}
