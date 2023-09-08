package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/subscription/service"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type SubscriptionHandler struct {
	logger logging.Logger

	subscriptionServ service.SubscriptionService
}

func NewSubscriptionHandler(logger logging.Logger, subscriptionServ service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger:           logger.With("handler", "subscription", "handlerType", "http"),
		subscriptionServ: subscriptionServ,
	}
}

func (h *SubscriptionHandler) HandleSubscribe() echo.HandlerFunc {
	type Request struct {
		ToUserNickname string `param:"nickname"`
	}

	return func(c echo.Context) error {
		const op = "SubscriptionHandler.HandleSubscribe"
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

		authPayload, _ := jwt.Extract(ctx)

		err = h.subscriptionServ.Subscribe(ctx, service.SubscribeDTO{
			FromUserID:     authPayload.UserID,
			ToUserNickname: req.ToUserNickname,
		})
		if err != nil {
			logger.Error("failed to subscribe", "error", err)

			if errors.Is(err, usermodel.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, usermodel.ErrUserNotFound.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusCreated)
	}
}
