package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/subscription/model"
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
		FindAllSubscriptionsByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]model.Subscription, error)
		CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (uuid.UUID, error)
		DeleteSubscription(ctx context.Context, dto DeleteSubscriptionDTO) error
		CountSubscriptionsByFromUserID(ctx context.Context, fromUserID uuid.UUID) (uint64, error)
		CountSubscriptionsByToUserID(ctx context.Context, toUserID uuid.UUID) (uint64, error)
	}
)
