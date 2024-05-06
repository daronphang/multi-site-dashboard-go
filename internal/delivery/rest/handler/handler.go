package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/delivery/sse"
	ws "multi-site-dashboard-go/internal/delivery/websocket"
	uc "multi-site-dashboard-go/internal/usecase"
)

var logger, _ = internal.WireLogger()

type RestHandler struct {
	UseCase *uc.UseCaseService
}

func NewRestHandler(uc *uc.UseCaseService) *RestHandler {
	return &RestHandler{UseCase: uc}
}

func (h *RestHandler) Heartbeat(c echo.Context) error {
	return c.String(http.StatusOK, "Multi-site dashboard is alive")
}

func (h *RestHandler) SSE(c echo.Context) error {
	sse.SSEHandler(c.Response().Writer, c.Request())
	return nil
}

func (h *RestHandler) Websocket(c echo.Context) error {
	ws.WebsocketHandler(c.Response().Writer, c.Request())
	return nil
}
