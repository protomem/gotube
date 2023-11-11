package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

type (
	CreateVideoDTO struct {
		Title         string
		Description   *string
		ThumbnailPath string
		VideoPath     string
		AuthorID      uuid.UUID
		Public        *bool
	}

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
