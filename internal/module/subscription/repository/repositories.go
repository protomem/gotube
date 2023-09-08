package repository

import (
	"context"

	"github.com/google/uuid"
)

type CreateSubscriptionDTO struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
}

type (
	SubscriptionRepository interface {
		CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (uuid.UUID, error)
	}
)
