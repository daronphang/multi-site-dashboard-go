package v1

import (
	"net/http"

	"multi-site-dashboard-go/api"

	"github.com/labstack/echo/v4"
)

type PtRequest struct {
	Name string `json:"name" validate:"required"`
	ID int `json:"id" validate:"required"`
}

func RegisterProdTrackingRoutes(g *echo.Group) {
	g.POST("/salesorder/status", func(c echo.Context) error {
		p := new(PtRequest)
		if err := api.ValidatePayload(c, p); err != nil {
			return api.NewHTTPValidationError(c, http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, p)
	})
}