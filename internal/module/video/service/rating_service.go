package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
)

var _ RatingService = (*RatingServiceImpl)(nil)

type LikeDTO struct {
	UserID  uuid.UUID
	VideoID uuid.UUID
}

type DislikeDTO struct {
	UserID  uuid.UUID
	VideoID uuid.UUID
}

type (
	RatingService interface {
		Like(ctx context.Context, dto LikeDTO) error
		Dislike(ctx context.Context, dto DislikeDTO) error
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

func (s *RatingServiceImpl) Like(ctx context.Context, dto LikeDTO) error {
	const op = "RatingService.Like"

	_, err := s.ratingRepo.CreateRating(ctx, repository.CreateRatingDTO{
		UserID: dto.UserID,
		VideID: dto.VideoID,
		Type:   model.Like,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RatingServiceImpl) Dislike(ctx context.Context, dto DislikeDTO) error {
	const op = "RatingService.Dislike"

	_, err := s.ratingRepo.CreateRating(ctx, repository.CreateRatingDTO{
		UserID: dto.UserID,
		VideID: dto.VideoID,
		Type:   model.Dislike,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
