package middleware

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func CustomRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI: true,
		LogLatency: true,
		BeforeNextFunc: func(c echo.Context) {
			// Request
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = io.ReadAll(c.Request().Body)
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			c.Set("payload", string(reqBody[:]))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info(
				"request logging",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("latency", v.Latency.String()),
				zap.String("payload", c.Get("payload").(string)),
				zap.Int64("bodySize", v.ResponseSize),
			)
			return nil
		},
	})
}


