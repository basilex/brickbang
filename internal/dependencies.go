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

type Container struct {
	DBPool *pgxpool.Pool

	AuxController  controller.IAuxController
	AuthController controller.IAuthController
}

func InitDependencies() *Container {
	cfg := config.Get()
	ctx := context.Background()

	// ===== DB Pool =====
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	// ===== Queries =====
	queries := dbs.New(dbPool)

	// ===== Repository layer =====
	auxRepo := repository.NewAuxRepository()
	authRepo := repository.NewAuthRepository(queries)

	// ===== Service layer =====
	auxService := service.NewAuxService(auxRepo)
	authService := service.NewAuthService(authRepo, cfg.JWTSecret, cfg.JWTAccessExpiration)

	// ===== Controller layer =====
	auxController := controller.NewAuxController(auxService)
	authController := controller.NewAuthController(authService)

	return &Container{
		DBPool:         dbPool,
		AuxController:  auxController,
		AuthController: authController,
	}
}
