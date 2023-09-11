package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/video/model"
)

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
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (uuid.UUID, error)
	}
)
