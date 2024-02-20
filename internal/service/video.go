package service

import "github.com/protomem/gotube/internal/repository"

var _ Video = (*VideoImpl)(nil)

type (
	Video interface{}

	VideoImpl struct {
		repo repository.Video
	}
)

func NewVideo(repo repository.Video) Video {
	return &VideoImpl{
		repo: repo,
	}
}
