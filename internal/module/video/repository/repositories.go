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

type (
	VideoRepository interface {
		FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error)
		FindAllVideosWherePublicAndSortByNew(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllVideosWherePublicAndSortByPopular(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllVideosByUserIDsAndWherePublicAndSortByNew(ctx context.Context, userIDs []uuid.UUID, opts FindAllVideosOptions) ([]model.Video, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)
	}
)
