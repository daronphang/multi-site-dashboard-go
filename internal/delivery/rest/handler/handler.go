package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"multi-site-dashboard-go/internal"
	uc "multi-site-dashboard-go/internal/usecase"
)

var logger, _ = internal.WireLogger()

type Handler struct {
	UseCase *uc.UseCaseService
}

func NewHandler(uc *uc.UseCaseService) *Handler {
	return &Handler{UseCase: uc}
}

func (h *Handler) Heartbeat(c echo.Context) error {
	return c.String(http.StatusOK, "Multi-site dashboard is alive")
}
