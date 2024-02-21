package repository

import (
	"context"

	"github.com/protomem/gotube/internal/model"
)

type (
	CreateVideoDTO struct {
		Title         string
		Description   string
		ThumbnailPath string
		VideoPath     string
		AuthorID      model.ID
		Public        bool
	}
)

type Video interface {
	Get(ctx context.Context, id model.ID) (model.Video, error)
	Create(ctx context.Context, dto CreateVideoDTO) (model.ID, error)
}
