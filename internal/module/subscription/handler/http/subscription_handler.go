package http

import (
	"github.com/protomem/gotube/internal/module/subscription/service"
	"github.com/protomem/gotube/pkg/logging"
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
