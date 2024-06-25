package rest

import (
	"context"
	"multi-site-dashboard-go/internal/delivery/rest/api/v1"
	rh "multi-site-dashboard-go/internal/delivery/rest/handler"
	cm "multi-site-dashboard-go/internal/delivery/rest/middleware"
	"net/http"

	uc "multi-site-dashboard-go/internal/usecase"
	cv "multi-site-dashboard-go/internal/validator"

	_ "multi-site-dashboard-go/docs"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"

	_ "multi-site-dashboard-go/statik"
)

type Server struct {
	Echo *echo.Echo
}

func NewServer(ctx context.Context, logger *zap.Logger, uc *uc.UseCaseService) (*Server, error) {
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

	// Register Swagger.
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/api/v1/swagger/", staticServer)
	eh := echo.WrapHandler(sh)
	baseGroup.GET("/swagger/*", eh)

	// Register custom handlers.
	e.Validator = cv.ProvideValidator()

	return &Server{
		Echo: e,
	}, nil
}