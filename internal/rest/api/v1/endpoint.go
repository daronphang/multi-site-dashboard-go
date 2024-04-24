package api

import (
	"github.com/labstack/echo/v4"
)

func (h *handler) RegisterBaseRoutes(g *echo.Group) {
	g.GET("/heartbeat", h.Heartbeat)
}

func (h *handler) RegisterPTRoutes(g *echo.Group) {
	g.POST("/salesorder", h.CreateSalesOrderPT)
}

func (h *handler) RegisterRTRoutes(g *echo.Group) {
	g.GET("/machines/:machine", h.GetMachineResourceUsageRT)
	g.POST("/machine", h.CreateMachineResourceUsageRT)
}