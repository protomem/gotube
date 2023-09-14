package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
)

var _ VideoService = (*VideoServiceImpl)(nil)

type CreateVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	UserID        uuid.UUID
}

type (
	VideoService interface {
		FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
	}

	VideoServiceImpl struct {
		videoRepo repository.VideoRepository
	}
)

func NewVideoService(videoRepo repository.VideoRepository) *VideoServiceImpl {
	return &VideoServiceImpl{
		videoRepo: videoRepo,
	}
}

func (s *VideoServiceImpl) FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "VideoService.FindOneVideo"

	video, err := s.videoRepo.FindOneVideo(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoServiceImpl) CreateVideo(ctx context.Context, dto CreateVideoDTO) (model.Video, error) {
	const op = "VideoService.CreateVideo"
	var err error

	// TODO: validate

	videoID, err := s.videoRepo.CreateVideo(ctx, repository.CreateVideoDTO(dto))
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	video, err := s.videoRepo.FindOneVideo(ctx, videoID)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}
