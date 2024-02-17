package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Subscription = (*Subscription)(nil)

type Subscription struct {
	logger logging.Logger
	db     database.DB
}

func NewSubscription(logger logging.Logger, db database.DB) *Subscription {
	return &Subscription{
		logger: logger.With("repository", "sqlite/subscription"),
		db:     db,
	}
}
