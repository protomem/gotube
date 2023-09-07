package http

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/auth/service"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type AuthMiddleware struct {
	logger logging.Logger

	authServ service.AuthService
}

func NewAuthMiddleware(logger logging.Logger, authServ service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger.With("middleware", "auth", "middlewareType", "http"),
		authServ: authServ,
	}
}

func (m *AuthMiddleware) Authenticator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "AuthMiddleware.Authenticator"

			ctx := c.Request().Context()
			logger := m.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				logger.Debug("no auth header")

				return next(c)
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) < 2 || headerParts[0] != "Bearer" {
				logger.Error("invalid auth header")

				return echo.ErrUnauthorized
			}

			accessToken := headerParts[1]

			_, authPayload, err := m.authServ.VerifyToken(ctx, accessToken)
			if err != nil {
				logger.Error("failed to verify tokens", "error", err)

				if errors.Is(err, usermodel.ErrUserNotFound) {
					return echo.ErrUnauthorized
				}

				return echo.ErrInternalServerError
			}

			ctx = jwt.Inject(ctx, authPayload)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (m *AuthMiddleware) Authorizer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "AuthMiddleware.Authorizer"

			ctx := c.Request().Context()
			logger := m.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			_, ok := jwt.Extract(ctx)
			if !ok {
				logger.Error("no auth payload")

				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
