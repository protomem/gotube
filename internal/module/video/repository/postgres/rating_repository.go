package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.RatingRepository = (*RatingRepository)(nil)

type RatingRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewRatingRepository(logger logging.Logger, db *database.DB) *RatingRepository {
	return &RatingRepository{
		logger:  logger.With("repository", "rating", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}
