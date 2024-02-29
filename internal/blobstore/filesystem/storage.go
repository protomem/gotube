package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/protomem/gotube/internal/blobstore"
	"github.com/protomem/gotube/pkg/logging"
)

var _ blobstore.Storage = (*Storage)(nil)

type Storage struct {
	logger     logging.Logger
	baseFolder string
}

func New(logger logging.Logger, folder string) (*Storage, error) {
	const op = "blobstore.New"

	baseFolder, err := filepath.Abs(folder)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{
		logger:     logger.With("component", "filesystem/blobstore"),
		baseFolder: baseFolder,
	}

	if err := s.initFolder(s.baseFolder); err != nil {
		return nil, fmt.Errorf("%s: init base folder: %w", op, err)
	}

	return s, nil
}

func (s *Storage) Get(ctx context.Context, folder, filename string) (blobstore.Object, error) {
	const op = "blobstore.Get"

	folder = s.fmtFolder(folder)
	filename = s.fmtFilename(folder, filename)

	s.logger.Debug("get object", "filename", filename)

	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return blobstore.Object{}, fmt.Errorf("%s: %w", op, blobstore.ErrObjectNotFound)
		}

		return blobstore.Object{}, fmt.Errorf("%s: %w", op, err)
	}

	return blobstore.Object{
		Type: s.resolveType(filename),
		Size: int64(len(file)),
		Body: io.NopCloser(bytes.NewBuffer(file)),
	}, nil
}

func (s *Storage) Put(ctx context.Context, folder, filename string, obj blobstore.Object) error {
	const op = "blobstore.Put"

	folder = s.fmtFolder(folder)
	filename = s.fmtFilename(folder, filename)

	s.logger.Debug("put object", "filename", filename)

	if err := s.initFolder(folder); err != nil {
		return fmt.Errorf("%s: init folder: %w", op, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := io.CopyN(file, obj.Body, obj.Size); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Del(ctx context.Context, folder, filename string) error {
	const op = "blobstore.Del"

	folder = s.fmtFolder(folder)
	filename = s.fmtFilename(folder, filename)

	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (*Storage) initFolder(folder string) error {
	return os.MkdirAll(folder, 0700)
}

func (s *Storage) fmtFolder(folder string) string {
	return filepath.Join(s.baseFolder, folder)
}

func (*Storage) fmtFilename(folder, filename string) string {
	return filepath.Join(folder, filename)
}

func (*Storage) resolveType(filename string) string {
	ext := filepath.Ext(filename)
	switch strings.TrimPrefix(strings.ToLower(ext), ".") {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
