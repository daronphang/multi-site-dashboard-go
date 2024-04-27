package rest

import (
	"context"
	"errors"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/repository"
	api "multi-site-dashboard-go/internal/rest/api/v1"
	cm "multi-site-dashboard-go/internal/rest/middleware"
	cv "multi-site-dashboard-go/internal/rest/validator"
	uc "multi-site-dashboard-go/internal/usecase"
	"os"
	"os/signal"
	"time"

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

	// Create repositories.
	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		return nil, err
	}
	repo := repository.New(db)

	// Create handler.
	svc := uc.NewUseCaseService(repo)
	h := api.NewHandler(svc)

	// Register routes.
	baseGroup := e.Group("/api/v1")
	h.RegisterBaseRoutes(baseGroup)

	ptGroup := baseGroup.Group("/pt")
	h.RegisterPTRoutes(ptGroup)

	rtGroup := baseGroup.Group("/rt")
	h.RegisterRTRoutes(rtGroup)

	// Register custom handlers.
	v := validator.New()
	err = v.RegisterValidation("int_required", cv.IntRequired)
	if err != nil {
		fmt.Printf("error registering custom validation: %v", err.Error())
	}
	e.Validator = &cv.CustomValidator{Validator: v}

	return &Server{Echo: e, Cfg: conf, Logger: logger, DB: db}, nil
}

func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		fmt.Printf("starting ECHO server in port %v", s.Cfg.Port)
		if err := s.Echo.Start(fmt.Sprintf(":%v", s.Cfg.Port)); err != nil {
			s.Logger.Fatal(fmt.Sprintf("failed to start server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Printf("performing graceful shutdown with timeout of %v...", 10*time.Second)
	s.DB.Close()
	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err.Error())
	}
}