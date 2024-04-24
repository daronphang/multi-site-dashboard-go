package rest

import (
	"context"
	"errors"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	api "multi-site-dashboard-go/internal/rest/api/v1"
	cm "multi-site-dashboard-go/internal/rest/middleware"
	cv "multi-site-dashboard-go/internal/rest/validator"
	uc "multi-site-dashboard-go/internal/usecase"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Echo *echo.Echo
	Logger *zap.Logger
	Cfg *config.Config
	DB *pgxpool.Pool
}

func NewServer() (*Server, error) {
	ctx := context.Background()
	wd, _ := os.Getwd()

	// Initialize dependencies.
	conf, err := config.ProvideConfig()
	if err != nil {
		msg := fmt.Sprintf("error reading config file: %v", err)
		return nil, errors.New(msg)
    }

	logger, err := internal.WireLogger()
	if err != nil {
		msg := fmt.Sprintf("error setting up logger: %v", err)
		return nil, errors.New(msg)
    }

	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		msg := fmt.Sprintf("error connecting to DB: %v", err)
		return nil, errors.New(msg)
	}

	// Migrate db.
	m, err := internal.WirePgMigrateInstance(wd)
	if err != nil {
		msg := fmt.Sprintf("error creating DB migration instance: %v", err)
		return nil, errors.New(msg)
	}
	if err := m.Up(); err.Error() != "no change" {
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
	svc := uc.NewUseCaseService()
	h := api.NewHandler(svc)

	// Register routes.
	baseGroup := e.Group("/api/v1")
	h.RegisterBaseRoutes(baseGroup)

	ptGroup := baseGroup.Group("/pt")
	h.RegisterPTRoutes(ptGroup)

	rtGroup := baseGroup.Group("/rt")
	h.RegisterRTRoutes(rtGroup)

	// Register custom validator.
	e.Validator = &cv.CustomValidator{Validator: validator.New()}

	return &Server{Echo: e, Cfg: conf, Logger: logger, DB: db}, nil
}

func (s *Server) Run() error {
	fmt.Printf("Starting ECHO server in port %v", s.Cfg.Port)
	if err := s.Echo.Start(fmt.Sprintf(":%v", s.Cfg.Port)); err != nil {
		msg := fmt.Sprintf("failed to start server: %v", err)
		return errors.New(msg)
	}
	return nil
}