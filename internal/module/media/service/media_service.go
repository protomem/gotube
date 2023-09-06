package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/storage"
)

var _ MediaService = (*MediaServiceImpl)(nil)

type SaveFileDTO struct {
	Folder string
	File   string
	Obj    storage.Object
}

type (
	MediaService interface {
		GetFile(ctx context.Context, folder, file string) (storage.Object, error)
		SaveFile(ctx context.Context, dto SaveFileDTO) error
	}

	MediaServiceImpl struct {
		bstore storage.BlobStorage
	}
)

func NewMediaService(bstore storage.BlobStorage) *MediaServiceImpl {
	return &MediaServiceImpl{
		bstore: bstore,
	}
}

func (s *MediaServiceImpl) GetFile(ctx context.Context, folder, file string) (storage.Object, error) {
	const op = "MediaService.GetFile"

	obj, err := s.bstore.GetObject(ctx, folder, file)
	if err != nil {
		return storage.Object{}, fmt.Errorf("%s: %w", op, err)
	}

	return obj, nil
}

func (s *MediaServiceImpl) SaveFile(ctx context.Context, dto SaveFileDTO) error {
	const op = "MediaService.SaveFile"

	err := s.bstore.PutObject(ctx, dto.Folder, dto.File, dto.Obj)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
