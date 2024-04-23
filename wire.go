//go:build wireinject

package main

import (
	"context"
	config "multi-site-dashboard-go/config"
	"multi-site-dashboard-go/repository"

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
	wire.Build(config.ProvideConfig, repository.ProvidePgConnPool)
	return &pgxpool.Pool{}, nil
}

func WirePgMigrateInstance(relativePath string) (*migrate.Migrate, error) {
	wire.Build(
		config.ProvideConfig,
		repository.ProvidePgDriver,
		repository.ProvidePgMigrateInstance,
	)
	return &migrate.Migrate{}, nil
}