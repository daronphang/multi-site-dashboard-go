package handler

import (
	"multi-site-dashboard-go/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *RestHandler) CreateSalesOrderPT(c echo.Context) error {
	c.Logger()
	p := new(domain.PTSalesOrder)
	if err := bindAndValidateRequestBody(c, p); err != nil {
		logger.Error(err.Error())
		return newHTTPValidationError(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}
