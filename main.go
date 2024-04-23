package main

import (
	"fmt"
	"multi-site-dashboard-go/api"
	v1 "multi-site-dashboard-go/api/v1"
	"multi-site-dashboard-go/config"
	customMiddleware "multi-site-dashboard-go/middleware"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize providers.
	cfg, err := config.ProvideConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v", err))
    }
	logger, err := WireLogger()
	if err != nil {
		panic(fmt.Sprintf("error setting up logger: %v", err))
    }

	// Migrate db.
	migrations := [1]string{
		"repository/resource_tracking/migrations",
	}
	for _, rp := range migrations {
		m, err := WirePgMigrateInstance(rp)
		if err != nil {
			panic(fmt.Sprintf("error creating PG migration instance: %v", err))
		}
		if err := m.Up(); err.Error() != "no change" {
			panic(fmt.Sprintf("error migrating db: %v", err))
		}
	}

	// Create app.
	e := echo.New()

	// Register middlewares.
	e.Use(
		middleware.CORS(),
		customMiddleware.CustomRequestLogger(logger),
	)

	// Register routes.
	defaultGroup := e.Group("/api/v1")
	v1.RegisterDefaultRoutes(defaultGroup)

	ptGroup := e.Group("/api/v1/pt")
	v1.RegisterProductionTrackingRoutes(ptGroup)

	rtGroup := e.Group("/api/v1/rt")
	v1.RegisterResourceTrackingRoutes(rtGroup)

	// Register custom validator.
	e.Validator = &api.CustomValidator{Validator: validator.New()}

	logger.Info(fmt.Sprintf("Starting ECHO server in port %v", cfg.Port))
	if err := e.Start(fmt.Sprintf(":%v", cfg.Port)); err != nil {
		msg := fmt.Sprintf("failed to start server: %v", err)
		logger.Error(msg)
		panic(msg)
	}
}