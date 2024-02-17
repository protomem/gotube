package sqlite

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

func New(logger logging.Logger, db database.DB) *repository.Repositories {
	return &repository.Repositories{
		User:         NewUser(logger, db),
		Subscription: NewSubscription(logger, db),
	}
}
