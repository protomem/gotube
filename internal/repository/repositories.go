package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
)

type (
	CreateUserDTO struct {
		Nickname string
		Password string
		Email    string
	}

	UpdateUserDTO struct {
		Nickname    *string
		Password    *string
		Email       *string
		Verified    *bool
		AvatarPath  *string
		Description *string
	}

	User interface {
		Get(ctx context.Context, id uuid.UUID) (model.User, error)
		GetByNickname(ctx context.Context, nickname string) (model.User, error)
		GetByEmail(ctx context.Context, email string) (model.User, error)

		Create(ctx context.Context, dto CreateUserDTO) (uuid.UUID, error)

		UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) error

		DeleteByNickname(ctx context.Context, nickname string) error
	}
)

type (
	CreateSubscriptionDTO struct {
		FromUserID uuid.UUID
		ToUserID   uuid.UUID
	}

	Subscription interface {
		FindByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]model.Subscription, error)

		GetByFromUserAndToUser(ctx context.Context, fromUserID, toUserID uuid.UUID) (model.Subscription, error)

		Create(ctx context.Context, dto CreateSubscriptionDTO) (uuid.UUID, error)

		Delete(ctx context.Context, id uuid.UUID) error
	}
)

type (
	CreateVideoDTO struct {
		Title         string
		Description   string
		ThumbnailPath string
		VideoPath     string
		AuthorID      uuid.UUID
		Public        bool
	}

	Video interface {
		Get(ctx context.Context, id uuid.UUID) (model.Video, error)

		Create(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)
	}
)
