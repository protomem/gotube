package service

import "github.com/protomem/gotube/internal/module/video/repository"

var _ VideoService = (*VideoServiceImpl)(nil)

type (
	VideoService interface{}

	VideoServiceImpl struct {
		videoRepo repository.VideoRepository
	}
)

func NewVideoService(videoRepo repository.VideoRepository) *VideoServiceImpl {
	return &VideoServiceImpl{
		videoRepo: videoRepo,
	}
}
