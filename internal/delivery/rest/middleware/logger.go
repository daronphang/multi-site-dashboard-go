package middleware

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var methodsWithPayload = []string{"POST", "PUT", "PATCH"}

func CustomRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI: true,
		LogLatency: true,
		BeforeNextFunc: func(c echo.Context) {


			// Request
			reqBody := []byte{}
			if c.Request().Body != nil && slices.Contains(methodsWithPayload, c.Request().Method) { // Read
				reqBody, _ = io.ReadAll(c.Request().Body)
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			c.Set("payload", string(reqBody[:]))
			c.Set("bodySize", byteCountIEC(len(reqBody)))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info(
				"request logging",
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.String("latency", v.Latency.String()),
				zap.String("payload", c.Get("payload").(string)),
				zap.String("bodySize", c.Get("bodySize").(string)),
			)
			return nil
		},
	})
}

func byteCountIEC(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

