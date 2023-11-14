package postgres

import (
	"database/sql"

	"github.com/protomem/gotube/pkg/logging"
)

type RatingRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewRatingRepository(logger logging.Logger, db *sql.DB) *RatingRepository {
	return &RatingRepository{
		logger: logger.With("repository", "postgres", "repositoryType", "rating"),
		db:     db,
	}
}
