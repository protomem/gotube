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
		FindAllNewVideosByUserNickname(ctx context.Context, userNickname string) ([]model.Video, error)
		FindAllPublicNewVideos(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		FindAllPublicPopularVideos(ctx context.Context, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		FindAllPublicNewVideosFromSubscriptions(ctx context.Context, userID uuid.UUID, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		SearchVideos(ctx context.Context, query string, opts FindAllVideosOptions) ([]model.Video, uint64, error)
		CreateVideo(ctx context.Context, dto CreateVideoDTO) (model.Video, error)
	}

	VideoServiceImpl struct {
		userServ userserv.UserService
		subServ  subserv.SubscriptionService

		videoRepo repository.VideoRepository
	}
)

func NewVideoService(
	userServ userserv.UserService,
	subServ subserv.SubscriptionService,
	videoRepo repository.VideoRepository,
) *VideoServiceImpl {
	return &VideoServiceImpl{
		userServ:  userServ,
		subServ:   subServ,
		videoRepo: videoRepo,
	}
}

func (s *VideoServiceImpl) FindOneVideo(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "VideoService.FindOneVideo"

	video, err := s.videoRepo.FindOneVideo(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.videoRepo.IncrementVideoView(ctx, id)
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (s *VideoServiceImpl) FindAllNewVideosByUserNickname(
	ctx context.Context,
	userNickname string,
) ([]model.Video, error) {
	const op = "VideoService.FindAllNewVideosByUserID"
	var err error

	user, err := s.userServ.FindOneUserByNickname(ctx, userNickname)
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	videos, err := s.videoRepo.FindAllVideosByUserIDAndSortByNew(ctx, user.ID)
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (s *VideoServiceImpl) FindAllPublicNewVideos(
	ctx context.Context,
	opts FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoService.FindAllPublicNewVideos"

	videos, count, err := s.videoRepo.FindAllVideosWherePublicAndSortByNew(ctx, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	return videos, count, nil
}

func (s *VideoServiceImpl) FindAllPublicPopularVideos(
	ctx context.Context,
	opts FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoService.FindAllPublicPopularVideos"

	videos, count, err := s.videoRepo.FindAllVideosWherePublicAndSortByPopular(ctx, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	return videos, count, nil
}

func (s *VideoServiceImpl) FindAllPublicNewVideosFromSubscriptions(
	ctx context.Context,
	userID uuid.UUID,
	opts FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoService.FindAllPublicNewVideosFromSubscriptions"

	user, err := s.userServ.FindOneUser(ctx, userID)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	subs, err := s.subServ.FindAllSubscriptionsByFromUserNickname(ctx, user.Nickname)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	userIDs := make([]uuid.UUID, 0, len(subs))
	for _, sub := range subs {
		userIDs = append(userIDs, sub.ToUser.ID)
	}

	videos, count, err := s.videoRepo.
		FindAllVideosByUserIDsAndWherePublicAndSortByNew(ctx, userIDs, repository.FindAllVideosOptions(opts))
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	return videos, count, nil
}

func (s *VideoServiceImpl) SearchVideos(
	ctx context.Context,
	query string,
	opts FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoService.SearchVideos"

	videos, count, err := s.videoRepo.FindAllVideosLikeByTitleAndWherePublicAndSortByNew(
		ctx,
		query,
		repository.FindAllVideosOptions(opts),
	)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	return videos, count, nil
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
