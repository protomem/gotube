package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type SubscriptionHandler struct {
	logger logging.Logger
	serv   service.Subscription
}

func NewSubscriptionHandler(logger logging.Logger, serv service.Subscription) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger: logger.With("handler", "subscription", "handlerType", "http"),
		serv:   serv,
	}
}

func (handl *SubscriptionHandler) Subscribe() http.HandlerFunc {
	type Request struct {
		FromUserID uuid.UUID `json:"fromUserId"`
		ToUserID   uuid.UUID `json:"toUserId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.SubscriptionHandler.Subscribe"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to decode request",
			})

			return
		}

		logger.Debug("subscribe request", "req", req)

		err = handl.serv.Subscribe(ctx, service.SubscribeDTO{
			FromUserID: req.FromUserID,
			ToUserID:   req.ToUserID,
		})
		if err != nil {
			logger.Error("failed to subscribe", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to subscribe",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (handl *SubscriptionHandler) Unsubscribe() http.HandlerFunc {
	type Request struct {
		FromUserID uuid.UUID `json:"fromUserId"`
		ToUserID   uuid.UUID `json:"toUserId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.SubscriptionHandler.Unsubscribe"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to decode request",
			})

			return
		}

		err = handl.serv.Unsubscribe(ctx, service.UnsubscribeDTO{
			FromUserID: req.FromUserID,
			ToUserID:   req.ToUserID,
		})
		if err != nil {
			logger.Error("failed to unsubscribe", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to unsubscribe",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
