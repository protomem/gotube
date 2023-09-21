package service

import "github.com/protomem/gotube/internal/module/video/repository"

var _ RatingService = (*RatingServiceImpl)(nil)

type (
	RatingService interface {
	}

	RatingServiceImpl struct {
		ratingRepo repository.RatingRepository
	}
)

func NewRatingService(ratingRepo repository.RatingRepository) *RatingServiceImpl {
	return &RatingServiceImpl{
		ratingRepo: ratingRepo,
	}
}
