package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/subscription/model"
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

func (h *SubscriptionHandler) HandleGetAllSubscriptions() echo.HandlerFunc {
	type Request struct {
		FromUserNickname string `param:"nickname"`
	}

	type Response struct {
		Subscriptions []model.Subscription `json:"subscriptions"`
	}

	return func(c echo.Context) error {
		const op = "SubscriptionHandler.HandleGetAllSubscriptions"
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

		subscriptions, err := h.subscriptionServ.FindAllSubscriptionsByFromUserNickname(ctx, req.FromUserNickname)
		if err != nil {
			logger.Error("failed to find all subscriptions", "error", err)

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, Response{Subscriptions: subscriptions})
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

func (h *SubscriptionHandler) HandleUnsubscribe() echo.HandlerFunc {
	type Request struct {
		ToUserNickname string `param:"nickname"`
	}

	return func(c echo.Context) error {
		const op = "SubscriptionHandler.HandleUnsubscribe"
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

		err = h.subscriptionServ.Unsubscribe(ctx, service.UnsubscribeDTO{
			FromUserID:     authPayload.UserID,
			ToUserNickname: req.ToUserNickname,
		})
		if err != nil {
			logger.Error("failed to unsubscribe", "error", err)

			if errors.Is(err, usermodel.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, usermodel.ErrUserNotFound.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (h *SubscriptionHandler) HandleStatistics() echo.HandlerFunc {
	type Request struct {
		UserNickname string `param:"nickname"`
	}

	type Response struct {
		CountSubscriptions uint64 `json:"subscriptions"`
		CountSubscribers   uint64 `json:"subscribers"`
	}

	return func(c echo.Context) error {
		const op = "SubscriptionHandler.HandleCountStatistics"
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

		statsDTO, err := h.subscriptionServ.Statistics(ctx, req.UserNickname)
		if err != nil {
			logger.Error("failed to count statistics", "error", err)

			if errors.Is(err, usermodel.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, usermodel.ErrUserNotFound.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, Response{
			CountSubscriptions: statsDTO.CountSubscriptions,
			CountSubscribers:   statsDTO.CountSubscribers,
		})
	}
}
