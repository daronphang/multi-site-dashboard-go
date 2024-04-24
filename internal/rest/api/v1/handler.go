package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"multi-site-dashboard-go/internal/domain"
	cv "multi-site-dashboard-go/internal/rest/validator"
	uc "multi-site-dashboard-go/internal/usecase"
)

type handler struct {
	Service *uc.UseCaseService
}

func NewHandler(uc *uc.UseCaseService) *handler {
	return &handler{Service: uc}
}

func (h *handler) Heartbeat(c echo.Context) error {
	return c.String(http.StatusOK, "Multi-site dashboard is alive")
}

func (h *handler) CreateSalesOrderPT(c echo.Context) error {
	p := new(domain.PTSalesOrder)
	if err := cv.ValidatePayload(c, p); err != nil {
		return cv.NewHTTPValidationError(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}

func (h *handler) GetMachineResourceUsageRT(c echo.Context) error {
	machine := c.Param("machine")
	err := h.Service.GetTimeSeriesMachineResourceUsageRT(c.Request().Context(), machine)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "hello world")
}

func (h *handler) CreateMachineResourceUsageRT(c echo.Context) error {
	p := new(domain.MachineResource)
	if err := cv.ValidatePayload(c, p); err != nil {
		return cv.NewHTTPValidationError(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}