package service

import "github.com/protomem/gotube/internal/storage"

type (
	MediaService interface{}

	MediaServiceImpl struct {
		bstore storage.BlobStorage
	}
)

func NewMediaService(bstore storage.BlobStorage) *MediaServiceImpl {
	return &MediaServiceImpl{
		bstore: bstore,
	}
}
