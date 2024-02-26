package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Comment = (*Comment)(nil)

type Comment struct {
	logger logging.Logger
	db     database.DB
}

func NewComment(logger logging.Logger, db database.DB) *Comment {
	return &Comment{
		logger: logger.With("repository", "sqlite/comment"),
		db:     db,
	}
}
