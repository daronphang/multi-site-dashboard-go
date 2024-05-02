package handler

import (
	"multi-site-dashboard-go/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAggMachineResourceUsageRT(c echo.Context) error {
	machine := c.Param("machine")
	lookBackPeriod := c.QueryParam("lookBackPeriod") // '1 hour', '1 day', '23 hours'
	timeBucket := c.QueryParam("timeBucket") // '5 minutes', '1 hour', '1 day'
	ctx := c.Request().Context()

	p := &domain.GetAggMachineResourceUsageParams{Machine: machine, TimeBucket: timeBucket, LookBackPeriod: lookBackPeriod}
	if err := bindAndValidateRequestBody(c, p); err != nil {
		logger.Error(err.Error())
		return newHTTPValidationError(c, http.StatusBadRequest, err)
	}
	rv, err := h.UseCase.GetAggMachineResourceUsageRT(ctx, p)

	if err != nil {
		logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, rv)
}

func (h *Handler) CreateMachineResourceUsageRT(c echo.Context) error {
	p := new(domain.CreateMachineResourceUsageParams)
	if err := bindAndValidateRequestBody(c, p); err != nil {
		logger.Error(err.Error())
		return newHTTPValidationError(c, http.StatusBadRequest, err)
	}
	rv, err := h.UseCase.CreateMachineResourceUsage(c.Request().Context(), p)
	if err != nil {
		logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, rv)
}