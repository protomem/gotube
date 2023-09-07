package service

import "github.com/protomem/gotube/internal/module/subscription/repository"

var _ SubscriptionService = (*SubscriptionServiceImpl)(nil)

type (
	SubscriptionService interface{}

	SubscriptionServiceImpl struct {
		subscriptionRepo repository.SubscriptionRepository
	}
)

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{
		subscriptionRepo: subscriptionRepo,
	}
}
