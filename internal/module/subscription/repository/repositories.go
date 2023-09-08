package repository

import (
	"context"

	"github.com/google/uuid"
)

type CreateSubscriptionDTO struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
}

type DeleteSubscriptionDTO struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
}

type (
	SubscriptionRepository interface {
		CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (uuid.UUID, error)
		DeleteSubscription(ctx context.Context, dto DeleteSubscriptionDTO) error
	}
)
