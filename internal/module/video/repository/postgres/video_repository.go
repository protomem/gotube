package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.VideoRepository = (*VideoRepository)(nil)

type VideoRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewVideoRepository(logger logging.Logger, db *database.DB) *VideoRepository {
	return &VideoRepository{
		logger:  logger.With("repository", "video", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}
