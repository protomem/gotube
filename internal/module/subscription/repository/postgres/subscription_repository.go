package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/subscription/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewSubscriptionRepository(logger logging.Logger, db *database.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger:  logger.With("repository", "subscription", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}
