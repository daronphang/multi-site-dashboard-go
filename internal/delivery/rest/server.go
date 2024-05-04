package rest

import (
	"context"
	"errors"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/rest/api/v1"
	rh "multi-site-dashboard-go/internal/delivery/rest/handler"
	cm "multi-site-dashboard-go/internal/delivery/rest/middleware"

	uc "multi-site-dashboard-go/internal/usecase"
	cv "multi-site-dashboard-go/internal/validator"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Echo *echo.Echo
}

func NewServer(ctx context.Context, cfg *config.Config, logger *zap.Logger, uc *uc.UseCaseService) (*Server, error) {
	wd, _ := os.Getwd()

	// Migrate db.
	m, err := internal.WirePgMigrateInstance(wd)
	if err != nil {
		msg := fmt.Sprintf("error creating DB migration instance: %v", err)
		return nil, errors.New(msg)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		msg := fmt.Sprintf("error migrating db: %v", err)
		return nil, errors.New(msg)
	}

	// Create server.
	e := echo.New()

	// Register middlewares.
	e.Use(
		middleware.CORS(),
		cm.CustomRequestLogger(logger),
	)

	// Create handler.
	rh := rh.NewRestHandler(uc)

	// Register routes.
	baseGroup := e.Group("/api/v1")
	api.RegisterBaseRoutes(baseGroup, rh)

	ptGroup := baseGroup.Group("/pt")
	api.RegisterPTRoutes(ptGroup, rh)

	rtGroup := baseGroup.Group("/rt")
	api.RegisterRTRoutes(rtGroup, rh)

	// Register custom handlers.
	e.Validator = cv.ProvideValidator()

	return &Server{
		Echo: e,
	}, nil
}