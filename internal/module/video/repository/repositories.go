package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/video/model"
)

type FindAllVideosOptions struct {
	Limit  uint64
	Offset uint64
}

type CreateVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	UserID        uuid.UUID
}

type CreateRatingDTO struct {
	UserID uuid.UUID
	VideID uuid.UUID
	Type   model.RatingType
}

type (
	VideoRepository interface {
		FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error)
		FindAllVideosByUserIDAndSortByNew(ctx context.Context, userID uuid.UUID) ([]model.Video, error)
		FindAllVideosWherePublicAndSortByNew(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllVideosWherePublicAndSortByPopular(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllVideosByUserIDsAndWherePublicAndSortByNew(ctx context.Context, userIDs []uuid.UUID, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllVideosLikeByTitleAndWherePublicAndSortByNew(ctx context.Context, title string, opts FindAllVideosOptions) ([]model.Video, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)
		IncrementVideoView(ctx context.Context, id uuid.UUID) error
	}

	RatingRepository interface {
		CreateRating(ctx context.Context, dto CreateRatingDTO) (uuid.UUID, error)
	}
)
