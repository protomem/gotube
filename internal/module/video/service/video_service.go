package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	subserv "github.com/protomem/gotube/internal/module/subscription/service"
	userserv "github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
)

var _ VideoService = (*VideoServiceImpl)(nil)

type FindAllVideosOptions struct {
	Limit  uint64
	Offset uint64
}

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
		FindAllPublicNewVideos(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllPublicPopularVideos(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, error)
		FindAllPublicNewVideosFromSubscriptions(ctx context.Context, userID uuid.UUID, opts FindAllVideosOptions) ([]model.Video, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
	}

	VideoServiceImpl struct {
		userServ userserv.UserService
		subServ  subserv.SubscriptionService

		videoRepo repository.VideoRepository
	}
)

func NewVideoService(userServ userserv.UserService, subServ subserv.SubscriptionService, videoRepo repository.VideoRepository) *VideoServiceImpl {
	return &VideoServiceImpl{
		userServ:  userServ,
		subServ:   subServ,
		videoRepo: videoRepo,
	}
}

func (s *VideoServiceImpl) FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "VideoService.FindOneVideo"

	// TODO: Increment video.View

	video, err := s.videoRepo.FindOneVideo(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoServiceImpl) FindAllPublicNewVideos(
	ctx context.Context,
	opts FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoService.FindAllPublicNewVideos"

	videos, err := s.videoRepo.FindAllVideosWherePublicAndSortByNew(ctx, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (s *VideoServiceImpl) FindAllPublicPopularVideos(
	ctx context.Context,
	opts FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoService.FindAllPublicPopularVideos"

	videos, err := s.videoRepo.FindAllVideosWherePublicAndSortByPopular(ctx, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (s *VideoServiceImpl) FindAllPublicNewVideosFromSubscriptions(
	ctx context.Context,
	userID uuid.UUID,
	opts FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoService.FindAllPublicNewVideosFromSubscriptions"

	user, err := s.userServ.FindOneUser(ctx, userID)
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	subs, err := s.subServ.FindAllSubscriptionsByFromUserNickname(ctx, user.Nickname)
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	userIDs := make([]uuid.UUID, 0, len(subs))
	for _, sub := range subs {
		userIDs = append(userIDs, sub.ToUser.ID)
	}

	videos, err := s.videoRepo.
		FindAllVideosByUserIDsAndWherePublicAndSortByNew(ctx, userIDs, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
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
