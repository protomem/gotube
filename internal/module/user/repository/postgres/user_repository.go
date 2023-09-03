package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/user/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepository(logger logging.Logger, db *database.DB) *UserRepository {
	return &UserRepository{
		logger:  logger.With("repository", "user", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}
