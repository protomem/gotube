package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Rating = (*Rating)(nil)

type Rating struct {
	logger logging.Logger
	db     database.DB
}

func NewRating(logger logging.Logger, db database.DB) *Rating {
	return &Rating{
		logger: logger.With("repository", "sqlite/rating"),
		db:     db,
	}
}
