package service

import "github.com/protomem/gotube/internal/repository"

var _ Rating = (*RatingImpl)(nil)

type (
	Rating interface{}

	RatingImpl struct {
		repo repository.Rating
	}
)

func NewRating(repo repository.Rating) *RatingImpl {
	return &RatingImpl{
		repo: repo,
	}
}
