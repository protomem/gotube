package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/module/auth/model"
	"github.com/protomem/gotube/internal/module/auth/service"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type AuthHandler struct {
	logger logging.Logger

	authServ service.AuthService
}

func NewAuthHandler(logger logging.Logger, authServ service.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:   logger.With("handler", "auth", "handlerType", "http"),
		authServ: authServ,
	}
}

func (h *AuthHandler) HandleRegister() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		model.PairTokens
		User usermodel.User `json:"user"`
	}

	return func(c echo.Context) error {
		const op = "AuthHandler.HandleRegister"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			logger.Error("failed to bind request", "error", err)

			return echo.ErrBadRequest
		}

		user, tokens, err := h.authServ.Register(ctx, service.RegisterDTO(req))
		if err != nil {
			logger.Error("failed to register user", "error", err)

			if errors.Is(err, usermodel.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, usermodel.ErrUserAlreadyExists.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusCreated, Response{
			PairTokens: tokens,
			User:       user,
		})
	}
}

func (h *AuthHandler) HandleLogin() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		model.PairTokens
		User usermodel.User `json:"user"`
	}

	return func(c echo.Context) error {
		const op = "AuthHandler.HandleLogin"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			logger.Error("failed to bind request", "error", err)

			return echo.ErrBadRequest
		}

		user, tokens, err := h.authServ.Login(ctx, service.LoginDTO(req))
		if err != nil {
			logger.Error("failed to login user", "error", err)

			if errors.Is(err, usermodel.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, usermodel.ErrUserNotFound.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, Response{
			PairTokens: tokens,
			User:       user,
		})
	}
}

func (h *AuthHandler) HandleRefreshTokens() echo.HandlerFunc {
	type Request struct {
		RefreshToken string `json:"refreshToken"`
	}

	type Response struct {
		model.PairTokens
	}

	return func(c echo.Context) error {
		const op = "AuthHandler.HandleRefreshTokens"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			logger.Error("failed to bind request", "error", err)

			return echo.ErrBadRequest
		}

		tokens, err := h.authServ.RefreshTokens(ctx, req.RefreshToken)
		if err != nil {
			logger.Error("failed to refresh tokens", "error", err)

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, Response{
			PairTokens: tokens,
		})
	}
}
