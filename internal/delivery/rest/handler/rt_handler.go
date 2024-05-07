package handler

import (
	"multi-site-dashboard-go/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// swagger:route GET /machines/{machine} ResourceTracking GetAggMachineResourceUsage
// Group MachineResourceUsage time series by timeBucket within lookBackPeriod from today and return aggregated median values
// responses:
// 	200: body:[]AggMachineResourceUsage
// 	400: body:HTTPValidationError
func (h *RestHandler) GetAggMachineResourceUsageRT(c echo.Context) error {
	machine := c.Param("machine")
	lookBackPeriod := c.QueryParam("lookBackPeriod") 
	timeBucket := c.QueryParam("timeBucket") 
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

// swagger:operation POST /machine ResourceTracking CreateMachineResourceUsage
// Add a time series entry for MachineResourceUsage
// ---
// parameters:
// - name: MachineResourceUsageParam
//   in: body
//   schema:
//     $ref: "#/definitions/CreateMachineResourceUsage"
// responses:
//   "200":
//     description: MachineResourceUsage
//     schema:
//       $ref: '#/definitions/MachineResourceUsage'
//   "400":
//     description: HTTPValidationError
//     schema:
// 	     $ref: '#/definitions/HTTPValidationError'
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