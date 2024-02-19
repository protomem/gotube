package repository

import (
	"context"

	"github.com/protomem/gotube/internal/model"
)

type CreateSubscriptionDTO struct {
	FromUserID model.ID
	ToUserID   model.ID
}

type Subscription interface {
	GetByFromUserIDAndToUserID(ctx context.Context, fromUserID, toUserID model.ID) (model.Subscription, error)
	Create(ctx context.Context, dto CreateSubscriptionDTO) (model.ID, error)
	Delete(ctx context.Context, id model.ID) error
}
