package main

import (
	"fmt"
	"multi-site-dashboard-go/api"
	v1 "multi-site-dashboard-go/api/v1"
	customMiddleware "multi-site-dashboard-go/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize providers.
	config, err := WireConfig()
	if err != nil {
		msg := fmt.Sprintf("failed to read config file: %v", err)
		panic(msg)
    }
	logger, err := WireLogger()
	if err != nil {
		msg := fmt.Sprintf("failed to setup logger: %v", err)
		panic(msg)
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

	ptGroup := e.Group("/api/v1/prod-tracking")
	v1.RegisterProdTrackingRoutes(ptGroup)

	// Register custom validator.
	e.Validator = &api.CustomValidator{Validator: validator.New()}

	logger.Info(fmt.Sprintf("Starting ECHO server in port %v", config.Port))
	if err := e.Start(fmt.Sprintf(":%v", config.Port)); err != nil {
		msg := fmt.Sprintf("failed to start server: %v", err)
		logger.Error(msg)
		panic(msg)
	}
}