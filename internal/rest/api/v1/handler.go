package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/domain"
	cv "multi-site-dashboard-go/internal/rest/validator"
	uc "multi-site-dashboard-go/internal/usecase"
)

var logger *zap.Logger

func init() {
	logger, _ = internal.WireLogger()
}

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
	c.Logger()
	p := new(domain.PTSalesOrder)
	if err := cv.ValidatePayload(c, p); err != nil {
		logger.Error(err.Error())
		return cv.NewHTTPValidationError(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}

func (h *handler) GetAggMachineResourceUsageRT(c echo.Context) error {
	machine := c.Param("machine")
	lookBackPeriod := c.QueryParam("lookBackPeriod") // '1 hour', '1 day', '23 hours'
	timeBucket := c.QueryParam("timeBucket") // '5 minutes', '1 hour', '1 day'
	ctx := c.Request().Context()

	p := &domain.GetAggMachineResourceUsageParams{Machine: machine, TimeBucket: timeBucket, LookBackPeriod: lookBackPeriod}
	if err := cv.ValidatePayload(c, p); err != nil {
		logger.Error(err.Error())
		return cv.NewHTTPValidationError(c, http.StatusBadRequest, err)
	}
	rv, err := h.Service.GetAggMachineResourceUsageRT(ctx, p)

	if err != nil {
		logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, rv)
}

func (h *handler) CreateMachineResourceUsageRT(c echo.Context) error {
	p := new(domain.CreateMachineResourceUsageParams)
	if err := cv.ValidatePayload(c, p); err != nil {
		logger.Error(err.Error())
		return cv.NewHTTPValidationError(c, http.StatusBadRequest, err)
	}
	rv, err := h.Service.CreateMachineResourceUsage(c.Request().Context(), p)
	if err != nil {
		logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, rv)
}