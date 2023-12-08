package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/pkg/logging"
)

var _ User = (*UserImpl)(nil)

type (
	User interface{}

	UserImpl struct {
		logger logging.Logger
		pdb    *pgxpool.Pool
	}
)

func NewUser(logger logging.Logger, pdb *pgxpool.Pool) *UserImpl {
	return &UserImpl{
		logger: logger.With("repository", "user", "repositoryType", "postgres"),
		pdb:    pdb,
	}
}
