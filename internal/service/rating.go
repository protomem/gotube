package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Rating = (*RatingImpl)(nil)

type (
	RatingDTO struct {
		UserID  model.ID
		VideoID model.ID
	}
)

type (
	Rating interface {
		Like(ctx context.Context, dto RatingDTO) error
		Dislike(ctx context.Context, dto RatingDTO) error
		Delete(ctx context.Context, dto RatingDTO) error
	}

	RatingImpl struct {
		repo repository.Rating
	}
)

func NewRating(repo repository.Rating) *RatingImpl {
	return &RatingImpl{
		repo: repo,
	}
}

func (s *RatingImpl) Like(ctx context.Context, dto RatingDTO) error {
	const op = "service.Rating.Like"

	rating, err := s.repo.Get(ctx, repository.RatingDTO(dto))
	if err != nil && !errors.Is(err, model.ErrRatingNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !errors.Is(err, model.ErrRatingNotFound) && !rating.Like {
		if err := s.repo.Delete(ctx, repository.RatingDTO(dto)); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	repoDTO := repository.CreateRatingDTO{
		RatingDTO: repository.RatingDTO(dto),
		Like:      true,
	}

	if _, err := s.repo.Create(ctx, repoDTO); err != nil && !errors.Is(err, model.ErrRatingExists) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RatingImpl) Dislike(ctx context.Context, dto RatingDTO) error {
	const op = "service.Rating.Dislike"

	rating, err := s.repo.Get(ctx, repository.RatingDTO(dto))
	if err != nil && !errors.Is(err, model.ErrRatingNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !errors.Is(err, model.ErrRatingNotFound) && rating.Like {
		if err := s.repo.Delete(ctx, repository.RatingDTO(dto)); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	repoDTO := repository.CreateRatingDTO{
		RatingDTO: repository.RatingDTO(dto),
		Like:      false,
	}

	if _, err := s.repo.Create(ctx, repoDTO); err != nil && !errors.Is(err, model.ErrRatingExists) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RatingImpl) Delete(ctx context.Context, dto RatingDTO) error {
	const op = "service.Rating.Delete"

	if err := s.repo.Delete(ctx, repository.RatingDTO(dto)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
