package api

import (
	"multi-site-dashboard-go/internal/delivery/rest/handler"

	"github.com/labstack/echo/v4"
)

func RegisterBaseRoutes(g *echo.Group, h * handler.Handler) {
	g.GET("/heartbeat", h.Heartbeat)
}

func RegisterPTRoutes(g *echo.Group, h *handler.Handler) {
	g.POST("/salesorder", h.CreateSalesOrderPT)
}

func RegisterRTRoutes(g *echo.Group, h *handler.Handler) {
	g.GET("/machines/:machine", h.GetAggMachineResourceUsageRT)
	g.POST("/machine", h.CreateMachineResourceUsageRT)
}