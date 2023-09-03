package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type CommonMiddleware struct {
	logger logging.Logger
}

func NewCommonMiddleware(logger logging.Logger) *CommonMiddleware {
	return &CommonMiddleware{
		logger: logger.With("middleware", "common", "middlewareType", "http"),
	}
}

func (m *CommonMiddleware) RequestID() echo.MiddlewareFunc {
	return echo.WrapMiddleware(requestid.Middleware())
}

func (m *CommonMiddleware) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogProtocol:  true,
		LogMethod:    true,
		LogRoutePath: true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()

			m.logger.Info(
				"incoming request",
				requestid.LogKey, requestid.Extract(ctx),
				"protocol", v.Protocol,
				"method", v.Method,
				"path", v.RoutePath,
				"status", v.Status,
				"latency", v.Latency,
				"ip", v.RemoteIP,
			)

			return nil
		},
	})
}

func (m *CommonMiddleware) Recovery() echo.MiddlewareFunc {
	return middleware.Recover()
}

func (m *CommonMiddleware) CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	})
}
