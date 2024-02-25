package repository

import (
	"context"

	"github.com/protomem/gotube/internal/model"
)

type (
	RatingDTO struct {
		UserID  model.ID
		VideoID model.ID
	}

	CreateRatingDTO struct {
		RatingDTO
		Like bool
	}
)

type Rating interface {
	Get(ctx context.Context, dto RatingDTO) (model.Rating, error)
	Create(ctx context.Context, dto CreateRatingDTO) (model.ID, error)
	Delete(ctx context.Context, dto RatingDTO) error
}
