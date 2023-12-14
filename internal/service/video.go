package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/validation"
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
	return validation.Validate(func(v *validation.Validator) {
		v.CheckField(
			validation.MinRunes(dto.Title, 5) &&
				validation.MaxRunes(dto.Title, 100),
			"title", "invalid title",
		)
		v.CheckField(validation.MaxRunes(dto.Description, 500), "description", "invalid description")
		v.CheckField(
			dto.ThumbnailPath == "" ||
				(validation.MaxRunes(dto.ThumbnailPath, 300) &&
					validation.IsURL(dto.ThumbnailPath)),
			"thumbnailPath", "invalid thumbnailPath",
		)
		v.CheckField(
			dto.VideoPath == "" ||
				(validation.MaxRunes(dto.VideoPath, 300) &&
					validation.IsURL(dto.VideoPath)),
			"videoPath", "invalid videoPath",
		)
	})
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
