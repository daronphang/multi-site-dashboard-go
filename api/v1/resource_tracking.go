package v1

import (
	"net/http"

	"multi-site-dashboard-go/api"

	"github.com/labstack/echo/v4"
)

type MachineResource struct {
	Machine string `json:"name" validate:"required"`
	Metric1 int `json:"metric1" validate:"required"`
	Metric2 int `json:"metric2" validate:"required"`
	Metric3 int `json:"metric3" validate:"required"`
}

func RegisterResourceTrackingRoutes(g *echo.Group) {
	g.POST("/machine", func(c echo.Context) error {
		p := new(MachineResource)
		if err := api.ValidatePayload(c, p); err != nil {
			return api.NewHTTPValidationError(c, http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, p)
	})
}