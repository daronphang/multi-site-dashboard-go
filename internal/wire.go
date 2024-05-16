//go:build wireinject

package internal

import (
	"context"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func WireLogger() (*zap.Logger, error) {
	wire.Build(config.ProvideConfig, config.ProvideLogger)
	return &zap.Logger{}, nil
}

func WirePgConnPool(ctx context.Context) (*pgxpool.Pool, error) {
	wire.Build(config.ProvideConfig, database.ProvidePgConnPool)
	return &pgxpool.Pool{}, nil
}

func WirePgMigrateInstance() (*migrate.Migrate, error) {
	wire.Build(
		config.ProvideConfig,
		database.ProvidePgDriver,
		database.ProvidePgMigrateInstance,
	)
	return &migrate.Migrate{}, nil
}