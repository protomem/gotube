package postgres

import (
	"database/sql"

	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Subscription = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewSubscriptionRepository(logger logging.Logger, db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger: logger.With("repository", "subscription", "repositoryType", "postgres"),
		db:     db,
	}
}
