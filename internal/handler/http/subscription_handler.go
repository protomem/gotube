package http

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
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
