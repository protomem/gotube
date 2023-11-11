package postgres

import (
	"database/sql"

	"github.com/protomem/gotube/pkg/logging"
)

type VideoRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewVideoRepository(logger logging.Logger, db *sql.DB) *VideoRepository {
	return &VideoRepository{
		logger: logger.With("repository", "video", "repositoryType", "postgres"),
		db:     db,
	}
}
