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
	UserID  uuid.UUID
	VideoID uuid.UUID
	Type    model.RatingType
}

type DeleteRatingDTO struct {
	UserID  uuid.UUID
	VideoID uuid.UUID
}

type (
	VideoRepository interface {
		FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error)
		FindAllVideosByUserIDAndSortByNew(ctx context.Context, userID uuid.UUID) ([]model.Video, error)
		FindAllVideosWherePublicAndSortByNew(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		FindAllVideosWherePublicAndSortByPopular(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		FindAllVideosByUserIDsAndWherePublicAndSortByNew(ctx context.Context, userIDs []uuid.UUID, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		FindAllVideosLikeByTitleAndWherePublicAndSortByNew(ctx context.Context, title string, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)
		IncrementVideoView(ctx context.Context, id uuid.UUID) error
	}

	RatingRepository interface {
		FindAllRatingsByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Rating, error)
		CreateRating(ctx context.Context, dto CreateRatingDTO) (uuid.UUID, error)
		DeleteRating(ctx context.Context, dto DeleteRatingDTO) error
	}
)
