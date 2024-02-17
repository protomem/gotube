package handler

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Subscription struct {
	logger logging.Logger
	serv   service.Subscription
}

func NewSubscription(logger logging.Logger, serv service.Subscription) *Subscription {
	return &Subscription{
		logger: logger.With("handler", "subscription"),
		serv:   serv,
	}
}

func (h *Subscription) Count() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Subscription.Get"))
}

func (h *Subscription) Subscribe() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Subscription.Subscribe"))
}

func (h *Subscription) Unsubscribe() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Subscription.Unsubscribe"))
}

func (h *Subscription) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, model.ErrUserNotFound) {
			httplib.WriteJSON(w, http.StatusNotFound, httplib.JSON{"message": model.ErrUserNotFound.Error()})
			return
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
