package rest

import (
	"context"
	"errors"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/rest/api/v1"
	"multi-site-dashboard-go/internal/delivery/rest/handler"
	cm "multi-site-dashboard-go/internal/delivery/rest/middleware"
	"multi-site-dashboard-go/internal/delivery/stream"
	"multi-site-dashboard-go/internal/repository"
	uc "multi-site-dashboard-go/internal/usecase"
	cv "multi-site-dashboard-go/internal/validator"
	"os"
	"time"

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
	cfg, err := config.ProvideConfig()
	if err != nil {
		msg := fmt.Sprintf("error reading config file: %v", err)
		return nil, errors.New(msg)
    }

	logger, err := internal.WireLogger()
	if err != nil {
		msg := fmt.Sprintf("error setting up logger: %v", err)
		return nil, errors.New(msg)
    }

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

	// Create UseCase.
	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		return nil, err
	}
	repo := repository.New(db)
	sp := stream.New(cfg)

	uc := uc.NewUseCaseService(repo, sp)

	// Create handler and register routes.
	h := handler.NewHandler(uc)
	baseGroup := e.Group("/api/v1")
	api.RegisterBaseRoutes(baseGroup, h)

	ptGroup := baseGroup.Group("/pt")
	api.RegisterPTRoutes(ptGroup, h)

	rtGroup := baseGroup.Group("/rt")
	api.RegisterRTRoutes(rtGroup, h)

	// Register custom handlers.
	e.Validator = cv.ProvideValidator()

	return &Server{
		Echo: e,
		Cfg: cfg,
		Logger: logger,
		DB: db,
	}, nil
}

func (s *Server) GracefulShutdown(ctx context.Context) {
	fmt.Printf("performing graceful shutdown with timeout of %v...", 10*time.Second)

	s.DB.Close()

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Logger.Fatal(fmt.Sprintf("failed to shutdown echo server: %v", err.Error()))
	}
}