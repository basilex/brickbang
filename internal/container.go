package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"brickbang/internal/client"
	"brickbang/internal/config"
	"brickbang/internal/module"
	"brickbang/internal/repository/dbc"
	"brickbang/internal/repository/dbs"
	"brickbang/internal/service"
)

type Container struct {
	DBPool  *pgxpool.Pool
	DBCache *client.RedisClient

	Metadata map[string]string

	// Modules
	AuxModule      module.IAuxModule
	AuthModule     module.IAuthModule
	RoleModule     module.IRoleModule
	GrantModule    module.IGrantModule
	CountryModule  module.ICountryModule
	CurrencyModule module.ICurrencyModule
}

func NewContainer() *Container {
	cfg := config.Get()
	ctx := context.Background()

	// Redis
	dbCache, err := client.NewRedisClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to init redis", "error", err)
		os.Exit(1)
	}

	// Postgres
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		slog.Error("failed to init postgres", "error", err)
		os.Exit(1)
	}

	metadata := Metadata()
	queries := dbs.New(dbPool)

	// --- NEW: Blacklist Repository + Service ---
	blacklistRepo := dbc.NewBlacklistRepository(dbCache.Client)
	blacklistService := service.NewBlacklistService(blacklistRepo)

	// Modules
	auxModule := module.NewAuxModule(metadata)
	authModule := module.NewAuthModule(queries, blacklistService)
	roleModule := module.NewRoleModule(queries)
	grantModule := module.NewGrantModule(queries)
	countryModule := module.NewCountryModule(queries)
	currencyModule := module.NewCurrencyModule(queries)

	return &Container{
		DBPool:         dbPool,
		DBCache:        dbCache,
		Metadata:       metadata,
		AuxModule:      auxModule,
		AuthModule:     authModule,
		RoleModule:     roleModule,
		GrantModule:    grantModule,
		CountryModule:  countryModule,
		CurrencyModule: currencyModule,
	}
}
