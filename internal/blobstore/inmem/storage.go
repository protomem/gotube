package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/protomem/gotube/internal/blobstore"
	"github.com/protomem/gotube/pkg/logging"
)

var _ blobstore.Storage = (*Storage)(nil)

type Storage struct {
	logger logging.Logger

	mux   sync.RWMutex
	store map[string]blobstore.Object
}

func New(logger logging.Logger) (*Storage, error) {
	return &Storage{
		logger: logger.With("component", "in-memory/blobstore"),
		store:  make(map[string]blobstore.Object),
	}, nil
}

func (s *Storage) Get(ctx context.Context, parent, name string) (blobstore.Object, error) {
	const op = "blobstore.Get"

	s.mux.RLock()
	defer s.mux.RUnlock()

	s.logger.WithContext(ctx).Debug("get object", "parent", parent, "name", name)

	obj, ok := s.store[s.fmtKey(parent, name)]
	if !ok {
		return blobstore.Object{}, fmt.Errorf("%s: %w", op, blobstore.ErrObjectNotFound)
	}

	return obj, nil
}

func (s *Storage) Put(ctx context.Context, parent, name string, obj blobstore.Object) error {
	const op = "blobstore.Put"

	copyObj, err := obj.Clone()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.logger.WithContext(ctx).Debug("put object", "parent", parent, "name", name, "objSize", copyObj.Size, "objType", copyObj.Type)

	s.store[s.fmtKey(parent, name)] = copyObj

	return nil
}

func (s *Storage) Del(ctx context.Context, parent, name string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.logger.WithContext(ctx).Debug("delete object", "parent", parent, "name", name)

	delete(s.store, s.fmtKey(parent, name))

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.store = make(map[string]blobstore.Object)

	return nil
}

func (s *Storage) fmtKey(parent, name string) string {
	return fmt.Sprintf("%s/%s", parent, name)
}
