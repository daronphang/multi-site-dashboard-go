package handler

import (
	"multi-site-dashboard-go/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// swagger:route GET /machines/{machine} ResourceTracking GetAggMachineResourceUsage
// responses:
// 	200: []AggMachineResourceUsage
// 	400: HTTPValidationError
func (h *RestHandler) GetAggMachineResourceUsageRT(c echo.Context) error {
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

// swagger:route POST /machine ResourceTracking CreateMachineResourceUsage
// responses:
// 	200: MachineResourceUsage
// 	400: HTTPValidationError
func (h *RestHandler) CreateMachineResourceUsageRT(c echo.Context) error {
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