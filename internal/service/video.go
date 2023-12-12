package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Video = (*VideoImpl)(nil)

type CreateVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	AuthorID      model.ID
	Public        bool
}

func (dto CreateVideoDTO) Validate() error {
	return nil
}

type (
	Video interface {
		FindNew(ctx context.Context, opts FindOptions) ([]model.Video, error)
		FindPopular(ctx context.Context, opts FindOptions) ([]model.Video, error)
		Get(ctx context.Context, id model.ID) (model.Video, error)
		Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
		Delete(ctx context.Context, id model.ID) error
	}

	VideoImpl struct {
		repo repository.Video
	}
)

func NewVideo(repo repository.Video) *VideoImpl {
	return &VideoImpl{
		repo: repo,
	}
}

func (s *VideoImpl) FindNew(ctx context.Context, opts FindOptions) ([]model.Video, error) {
	const op = "service:Video.FindNew"

	videos, err := s.repo.Find(ctx, repository.FindOptions(opts))
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (s *VideoImpl) FindPopular(ctx context.Context, opts FindOptions) ([]model.Video, error) {
	const op = "service:Video.FindPopular"

	videos, err := s.repo.FindOrderByViews(ctx, repository.FindOptions(opts))
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (s *VideoImpl) Get(ctx context.Context, id model.ID) (model.Video, error) {
	const op = "service:Video.Get"

	video, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoImpl) Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error) {
	const op = "service:Video.Create"

	if err := dto.Validate(); err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	videoID, err := s.repo.Create(ctx, repository.CreateVideoDTO(dto))
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	video, err := s.repo.Get(ctx, videoID)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoImpl) Delete(ctx context.Context, id model.ID) error {
	const op = "service:Video.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
