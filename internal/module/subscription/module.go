package subscription

import (
	"github.com/protomem/gotube/internal/database"
	handlhttp "github.com/protomem/gotube/internal/module/subscription/handler/http"
	"github.com/protomem/gotube/internal/module/subscription/repository"
	repopostgres "github.com/protomem/gotube/internal/module/subscription/repository/postgres"
	"github.com/protomem/gotube/internal/module/subscription/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.SubscriptionHandler
	service.SubscriptionService
	repository.SubscriptionRepository
}

func New(logger logging.Logger, db *database.DB) *Module {
	logger = logger.With("module", "subscription")

	subscriptionRepo := repopostgres.NewSubscriptionRepository(logger, db)
	subscriptionServ := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandl := handlhttp.NewSubscriptionHandler(logger, subscriptionServ)

	return &Module{
		SubscriptionHandler:    subscriptionHandl,
		SubscriptionService:    subscriptionServ,
		SubscriptionRepository: subscriptionRepo,
	}
}
