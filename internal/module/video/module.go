package video

import (
	"github.com/protomem/gotube/internal/database"
	subserv "github.com/protomem/gotube/internal/module/subscription/service"
	userserv "github.com/protomem/gotube/internal/module/user/service"
	handlhttp "github.com/protomem/gotube/internal/module/video/handler/http"
	"github.com/protomem/gotube/internal/module/video/repository"
	repopostgres "github.com/protomem/gotube/internal/module/video/repository/postgres"
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.VideoHandler
	*handlhttp.RatingHandler

	service.VideoService
	service.RatingService

	repository.VideoRepository
	repository.RatingRepository
}

func New(
	logger logging.Logger,
	db *database.DB,
	userServ userserv.UserService,
	subServ subserv.SubscriptionService,
) *Module {
	logger = logger.With("module", "video")

	videoRepo := repopostgres.NewVideoRepository(logger, db)
	ratingRepo := repopostgres.NewRatingRepository(logger, db)

	videoServ := service.NewVideoService(userServ, subServ, videoRepo)
	ratingServ := service.NewRatingService(ratingRepo)

	videoHandl := handlhttp.NewVideoHandler(logger, videoServ)
	ratingHandl := handlhttp.NewRatingHandler(logger, ratingServ)

	return &Module{
		VideoHandler:     videoHandl,
		RatingHandler:    ratingHandl,
		VideoService:     videoServ,
		RatingService:    ratingServ,
		VideoRepository:  videoRepo,
		RatingRepository: ratingRepo,
	}
}
