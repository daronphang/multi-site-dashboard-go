// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go.uber.org/zap"
	"multi-site-dashboard-go/config"
)

// Injectors from wire.go:

func WireLogger() (*zap.Logger, error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return nil, err
	}
	logger, err := config.ProvideLogger(configConfig)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func WireConfig() (*config.Config, error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}