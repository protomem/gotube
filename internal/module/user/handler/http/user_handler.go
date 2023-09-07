package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type UserHandler struct {
	logger logging.Logger

	userServ service.UserService
}

func NewUserHandler(logger logging.Logger, userServ service.UserService) *UserHandler {
	return &UserHandler{
		logger:   logger.With("handler", "user", "handlerType", "http"),
		userServ: userServ,
	}
}

func (h *UserHandler) HandleGetUser() echo.HandlerFunc {
	type Request struct {
		Nickname string `param:"nickname"`
	}

	type Response struct {
		User model.User `json:"user"`
	}

	return func(c echo.Context) error {
		const op = "UserHandler.HandleGetUser"
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

		user, err := h.userServ.FindOneUserByNickname(ctx, req.Nickname)
		if err != nil {
			logger.Error("failed to find user", "error", err)

			if errors.Is(err, model.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, Response{User: user})
	}
}
