package video

import (
	"github.com/protomem/gotube/internal/database"
	handlhttp "github.com/protomem/gotube/internal/module/video/handler/http"
	"github.com/protomem/gotube/internal/module/video/repository"
	repopostgres "github.com/protomem/gotube/internal/module/video/repository/postgres"
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.VideoHandler
	service.VideoService
	repository.VideoRepository
}

func New(logger logging.Logger, db *database.DB) *Module {
	logger = logger.With("module", "video")

	videoRepo := repopostgres.NewVideoRepository(logger, db)
	videoServ := service.NewVideoService(videoRepo)
	videoHandl := handlhttp.NewVideoHandler(logger, videoServ)

	return &Module{
		VideoHandler:    videoHandl,
		VideoService:    videoServ,
		VideoRepository: videoRepo,
	}
}
