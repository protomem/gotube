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
	FindVideosOptions struct {
		Limit  uint64
		Offset uint64
	}

	CreateVideoDTO struct {
		Title         string
		Description   string
		ThumbnailPath string
		VideoPath     string
		AuthorID      uuid.UUID
		Public        bool
	}

	UpdateVideoDTO struct {
		Title         *string
		Description   *string
		ThumbnailPath *string
		VideoPath     *string
		Public        *bool
	}

	Video interface {
		FindAllPublicSortByCreatedAt(ctx context.Context, opts FindVideosOptions) ([]model.Video, error)
		FindAllPublicSortByViews(ctx context.Context, opts FindVideosOptions) ([]model.Video, error)
		FindByAuthorIDsSortByCreatedAt(ctx context.Context, authorIDs []uuid.UUID, opts FindVideosOptions) ([]model.Video, error)
		FindByAuthorNicknameSortByCreatedAt(ctx context.Context, nickname string) ([]model.Video, error)
		FindPublicByLikeTitleSortByCreatedAt(ctx context.Context, query string, opts FindVideosOptions) ([]model.Video, error)

		Get(ctx context.Context, id uuid.UUID) (model.Video, error)
		GetPublic(ctx context.Context, id uuid.UUID) (model.Video, error)

		Create(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)

		Update(ctx context.Context, id uuid.UUID, dto UpdateVideoDTO) error

		Delete(ctx context.Context, id uuid.UUID) error
	}
)

type (
	CreateRatingDTO struct {
		Like    bool
		VideoID uuid.UUID
		UserID  uuid.UUID
	}

	Rating interface {
		GetByVideoIDAndUserID(ctx context.Context, videoID, userID uuid.UUID) (model.Rating, error)

		Create(ctx context.Context, dto CreateRatingDTO) (uuid.UUID, error)

		Delete(ctx context.Context, id uuid.UUID) error
	}
)
