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

	Metadata map[string]string

	AuxController  controller.IAuxController
	KeyController  controller.IKeyController
	AuthController controller.IAuthController
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

	// Repository layer: responsible for database access
	auxRepo := repository.NewAuxRepository()
	keyRepo := repository.NewKeyRepository(queries)
	authRepo := repository.NewAuthRepository(queries)

	// Service layer: contains business logic
	auxService := service.NewAuxService(metadata, auxRepo)
	keyService := service.NewKeyService(keyRepo)
	authService := service.NewAuthService(authRepo, cfg.JWTSecret, cfg.JWTAccessExpiration)

	// Controller layer: handles HTTP requests/responses
	auxController := controller.NewAuxController(auxService)
	keyController := controller.NewKeyController(keyService)
	authController := controller.NewAuthController(authService)

	// Return a fully initialized dependency container
	return &Container{
		DBPool:         dbPool,
		Metadata:       metadata,
		AuxController:  auxController,
		KeyController:  keyController,
		AuthController: authController,
	}
}
