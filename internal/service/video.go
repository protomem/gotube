package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

type (
	FindVideosOptions struct {
		Limit  uint64
		Offset uint64
	}

	CreateVideoDTO struct {
		Title         string
		Description   *string
		ThumbnailPath string
		VideoPath     string
		AuthorID      uuid.UUID
		Public        *bool
	}

	Video interface {
		FindAllPublic(ctx context.Context, opts FindVideosOptions) ([]model.Video, error)

		Get(ctx context.Context, id uuid.UUID) (model.Video, error)
		GetPublic(ctx context.Context, id uuid.UUID) (model.Video, error)

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

func (serv *VideoImpl) FindAllPublic(ctx context.Context, opts FindVideosOptions) ([]model.Video, error) {
	const op = "service.Video.FindAllPublic"

	videos, err := serv.repo.FindAllPublic(ctx, repository.FindVideosOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (serv *VideoImpl) Get(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "service.Video.Get"

	video, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (serv *VideoImpl) GetPublic(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "service.Video.GetPublic"

	video, err := serv.repo.GetPublic(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (serv *VideoImpl) Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error) {
	const op = "service.Video.Create"
	var err error

	// TODO: Valiate ...

	createData := repository.CreateVideoDTO{
		Title:         dto.Title,
		ThumbnailPath: dto.ThumbnailPath,
		VideoPath:     dto.VideoPath,
		AuthorID:      dto.AuthorID,
	}

	if dto.Description != nil {
		createData.Description = *dto.Description
	}

	if dto.Public != nil {
		createData.Public = *dto.Public
	} else {
		createData.Public = true
	}

	id, err := serv.repo.Create(ctx, createData)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	video, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}
