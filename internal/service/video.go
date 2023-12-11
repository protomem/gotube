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
		Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
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
