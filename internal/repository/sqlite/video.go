package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Video = (*Video)(nil)

type Video struct {
	logger logging.Logger
	db     database.DB
}

func NewVideo(logger logging.Logger, db database.DB) *Video {
	return &Video{
		logger: logger.With("repository", "sqlite/video"),
		db:     db,
	}
}
