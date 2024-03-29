package repository

import (
	"context"

	"github.com/protomem/gotube/internal/model"
)

type (
	CreateCommentDTO struct {
		Message  string
		VideoID  model.ID
		AuthorID model.ID
	}
)

type Comment interface {
	FindByVideo(ctx context.Context, videoID model.ID, opts FindOptions) ([]model.Comment, error)
	Get(ctx context.Context, id model.ID) (model.Comment, error)
	Create(ctx context.Context, dto CreateCommentDTO) (model.ID, error)
	Delete(ctx context.Context, id model.ID) error
}
