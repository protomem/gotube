package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Rating = (*RatingImpl)(nil)

type (
	LikeDTO struct {
		VideoID uuid.UUID
		UserID  uuid.UUID
	}

	DislikeDTO struct {
		VideoID uuid.UUID
		UserID  uuid.UUID
	}

	Rating interface {
		FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Rating, error)

		Like(ctx context.Context, dto LikeDTO) error
		Dislike(ctx context.Context, dto DislikeDTO) error

		DeleteByVideoIDAndUserID(ctx context.Context, videoID, userID uuid.UUID) error
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

func (serv *RatingImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Rating, error) {
	const op = "service.Rating.FindByVideoID"

	ratings, err := serv.repo.FindByVideoID(ctx, videoID)
	if err != nil {
		return []model.Rating{}, fmt.Errorf("%s: %w", op, err)
	}

	return ratings, nil
}

func (serv *RatingImpl) Like(ctx context.Context, dto LikeDTO) error {
	const op = "service.Rating.Like"
	var err error

	rating, err := serv.repo.GetByVideoIDAndUserID(ctx, dto.VideoID, dto.UserID)
	if err != nil && !errors.Is(err, model.ErrRatingNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil && !rating.Like {
		err = serv.repo.Delete(ctx, rating.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	_, err = serv.repo.Create(ctx, repository.CreateRatingDTO{
		Like:    true,
		VideoID: dto.VideoID,
		UserID:  dto.UserID,
	})
	if err != nil && !errors.Is(err, model.ErrRatingExists) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (serv *RatingImpl) Dislike(ctx context.Context, dto DislikeDTO) error {
	const op = "service.Rating.Dislike"
	var err error

	rating, err := serv.repo.GetByVideoIDAndUserID(ctx, dto.VideoID, dto.UserID)
	if err != nil && !errors.Is(err, model.ErrRatingNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil && rating.Like {
		err = serv.repo.Delete(ctx, rating.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	_, err = serv.repo.Create(ctx, repository.CreateRatingDTO{
		Like:    false,
		VideoID: dto.VideoID,
		UserID:  dto.UserID,
	})
	if err != nil && !errors.Is(err, model.ErrRatingExists) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (serv *RatingImpl) DeleteByVideoIDAndUserID(ctx context.Context, videoID, userID uuid.UUID) error {
	const op = "service.Rating.DeleteByVideoIDAndUserID"
	var err error

	rating, err := serv.repo.GetByVideoIDAndUserID(ctx, videoID, userID)
	if err != nil {
		if errors.Is(err, model.ErrRatingNotFound) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	err = serv.repo.Delete(ctx, rating.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
