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
	GetByFromUserAndToUser(ctx context.Context, fromUserID, toUserID model.ID) (model.Subscription, error)
	CountByToUser(ctx context.Context, toUserID model.ID) (int64, error)
	Create(ctx context.Context, dto CreateSubscriptionDTO) (model.ID, error)
	Delete(ctx context.Context, id model.ID) error
}
