package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Video = (*VideoImpl)(nil)

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

	UpdateVideoDTO struct {
		Title         *string
		Description   *string
		ThumbnailPath *string
		VideoPath     *string
		Public        *bool
	}

	Video interface {
		FindNew(ctx context.Context, opts FindVideosOptions) ([]model.Video, error)
		FindPopular(ctx context.Context, opts FindVideosOptions) ([]model.Video, error)
		FindByAuthorNickname(ctx context.Context, nickname string) ([]model.Video, error)

		Get(ctx context.Context, id uuid.UUID) (model.Video, error)
		GetPublic(ctx context.Context, id uuid.UUID) (model.Video, error)

		Create(ctx context.Context, dto CreateVideoDTO) (model.Video, error)

		Update(ctx context.Context, id uuid.UUID, dto UpdateVideoDTO) (model.Video, error)

		Delete(ctx context.Context, id uuid.UUID) error
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

func (serv *VideoImpl) FindNew(ctx context.Context, opts FindVideosOptions) ([]model.Video, error) {
	const op = "service.Video.FindNew"

	videos, err := serv.repo.FindAllPublicSortByCreatedAt(ctx, repository.FindVideosOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (serv *VideoImpl) FindPopular(ctx context.Context, opts FindVideosOptions) ([]model.Video, error) {
	const op = "service.Video.FindPopular"

	videos, err := serv.repo.FindAllPublicSortByViews(ctx, repository.FindVideosOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (serv *VideoImpl) FindByAuthorNickname(ctx context.Context, nickname string) ([]model.Video, error) {
	const op = "service.Video.FindByAuthorNickname"

	// TODO: Valiate ...

	videos, err := serv.repo.FindByAuthorNicknameSortByCreatedAt(ctx, nickname)
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

func (serv *VideoImpl) Update(ctx context.Context, id uuid.UUID, dto UpdateVideoDTO) (model.Video, error) {
	const op = "service.Video.Update"
	var err error

	// TODO: Valiate ...

	_, err = serv.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	updateData := repository.UpdateVideoDTO{}

	if dto.Title != nil {
		updateData.Title = dto.Title
	}

	if dto.Description != nil {
		updateData.Description = dto.Description
	}

	if dto.ThumbnailPath != nil {
		updateData.ThumbnailPath = dto.ThumbnailPath
	}

	if dto.VideoPath != nil {
		updateData.VideoPath = dto.VideoPath
	}

	if dto.Public != nil {
		updateData.Public = dto.Public
	}

	err = serv.repo.Update(ctx, id, updateData)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	newVideo, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return newVideo, nil
}

func (serv *VideoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "service.Video.Delete"

	err := serv.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
