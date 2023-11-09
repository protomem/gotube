package service

import "github.com/protomem/gotube/internal/repository"

var _ Subscription = (*SubscriptionImpl)(nil)

type (
	Subscription interface {
	}

	SubscriptionImpl struct {
		repo repository.Subscription
	}
)

func NewSubscription(repo repository.Subscription) *SubscriptionImpl {
	return &SubscriptionImpl{
		repo: repo,
	}
}
