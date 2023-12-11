package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/pkg/logging"
)

type (
	Video interface{}

	VideoImpl struct {
		logger logging.Logger
		pdb    *pgxpool.Pool
	}
)

func NewVideo(logger logging.Logger, pdb *pgxpool.Pool) *UserImpl {
	return &UserImpl{
		logger: logger.With("repository", "video", "repositoryType", "postgres"),
		pdb:    pdb,
	}
}
