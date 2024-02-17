package service

import "github.com/protomem/gotube/internal/repository"

var _ Subscription = (*SubscriptionImpl)(nil)

type (
	Subscription interface {
	}

	SubscriptionImpl struct {
		repo     repository.Subscription
		userServ User
	}
)

func NewSubscription(repo repository.Subscription, userServ User) *SubscriptionImpl {
	return &SubscriptionImpl{
		repo:     repo,
		userServ: userServ,
	}
}
