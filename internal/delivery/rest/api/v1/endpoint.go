package api

import (
	"multi-site-dashboard-go/internal/delivery/rest/handler"

	"github.com/labstack/echo/v4"
)

func RegisterBaseRoutes(g *echo.Group, h * handler.RestHandler) {
	g.GET("/heartbeat", h.Heartbeat)
	g.GET("/sse", h.SSE)
	g.GET("/ws", h.Websocket)
}

func RegisterPTRoutes(g *echo.Group, h *handler.RestHandler) {
	g.POST("/salesorder", h.CreateSalesOrderPT)
}

func RegisterRTRoutes(g *echo.Group, h *handler.RestHandler) {
	g.GET("/machines/:machine", h.GetAggMachineResourceUsageRT)
	g.POST("/machine", h.CreateMachineResourceUsageRT)
}