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

	UpdateVideoDTO struct {
		Title         *string
		Description   *string
		ThumbnailPath *string
		VideoPath     *string
		Public        *bool
	}
)

type Video interface {
	FindSortByCreatedAtWherePublic(ctx context.Context, opts FindOptions) ([]model.Video, error)
	Get(ctx context.Context, id model.ID) (model.Video, error)
	Create(ctx context.Context, dto CreateVideoDTO) (model.ID, error)
	Update(ctx context.Context, id model.ID, dto UpdateVideoDTO) error
	Delete(ctx context.Context, id model.ID) error
}
