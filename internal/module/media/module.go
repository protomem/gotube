package media

import (
	handlhttp "github.com/protomem/gotube/internal/module/media/handler/http"
	"github.com/protomem/gotube/internal/module/media/service"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.MediaHandler
	service.MediaService
}

func New(logger logging.Logger, bstore storage.BlobStorage) *Module {
	logger = logger.With("module", "media")

	mediaServ := service.NewMediaService(bstore)
	mediaHandl := handlhttp.NewMediaHandler(logger, mediaServ)

	return &Module{
		MediaHandler: mediaHandl,
		MediaService: mediaServ,
	}
}
