package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.User = (*User)(nil)

type User struct {
	logger logging.Logger
	db     database.DB
}

func NewUser(logger logging.Logger, db database.DB) *User {
	return &User{
		logger: logger.With("repository", "sqlite/user"),
		db:     db,
	}
}
