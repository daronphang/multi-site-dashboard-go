//go:build wireinject

package main

import (
	config "multi-site-dashboard-go/config"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func WireLogger() (*zap.Logger, error) {
	wire.Build(config.ProvideConfig, config.ProvideLogger)
	return &zap.Logger{}, nil
}

func WireConfig() (*config.Config, error) {
	wire.Build(config.ProvideConfig)
	return &config.Config{}, nil
}