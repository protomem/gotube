package service

import (
	"context"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Video = (*VideoImpl)(nil)

type (
	CreateVideoDTO struct {
		Title         string
		Description   *string
		ThumbnailPath string
		VideoPath     string
		AuthorID      model.ID
		Public        *bool
	}

	UpdateVideoDTO struct {
		Title         *string
		Description   *string
		ThumbnailPath *string
		VideoPath     *string
		Public        *bool
	}
)

type (
	Video interface {
		Get(ctx context.Context, id model.ID) (model.Video, error)
		Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
		Update(ctx context.Context, id model.ID, dto UpdateVideoDTO) (model.Video, error)
		Delete(ctx context.Context, id model.ID) error
	}

	VideoImpl struct {
		repo repository.Video
	}
)

func NewVideo(repo repository.Video) Video {
	return &VideoImpl{
		repo: repo,
	}
}

func (s *VideoImpl) Get(ctx context.Context, id model.ID) (model.Video, error) {
	const op = "service.Video.Get"

	video, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoImpl) Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error) {
	const op = "service.Video.Create"

	// TODO: Add validation

	repoDTO := repository.CreateVideoDTO{
		Title:         dto.Title,
		Description:   s.autoGenerateVideoDescription(),
		ThumbnailPath: dto.ThumbnailPath,
		VideoPath:     dto.VideoPath,
		AuthorID:      dto.AuthorID,
		Public:        true,
	}
	if dto.Description != nil {
		repoDTO.Description = *dto.Description
	}
	if dto.Public != nil {
		repoDTO.Public = *dto.Public
	}

	id, err := s.repo.Create(ctx, repoDTO)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	video, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoImpl) Update(ctx context.Context, id model.ID, dto UpdateVideoDTO) (model.Video, error) {
	const op = "service.Video.Update"

	if _, err := s.repo.Get(ctx, id); err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.Update(ctx, id, repository.UpdateVideoDTO(dto)); err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	newVideo, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return newVideo, nil
}

func (s *VideoImpl) Delete(ctx context.Context, id model.ID) error {
	const op = "service.Video.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *VideoImpl) autoGenerateVideoDescription() string {
	return fmt.Sprintf("Auto generated description %d/%s", time.Now().Year(), time.Now().Month().String())
}
