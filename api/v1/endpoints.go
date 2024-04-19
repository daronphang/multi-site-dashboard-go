package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterDefaultRoutes(g *echo.Group) {
	g.GET("/heartbeat", func(c echo.Context) error {
		return c.String(http.StatusOK, "Multi-site dashboard is alive!")
	})
}