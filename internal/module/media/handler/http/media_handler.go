package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/module/media/service"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

const _maxFileSize = 100 * 1024 * 1024 // 100 MB

type MediaHandler struct {
	logger logging.Logger

	mediaServ service.MediaService
}

func NewMediaHandler(logger logging.Logger, mediaServ service.MediaService) *MediaHandler {
	return &MediaHandler{
		logger:    logger.With("handler", "media", "handlerType", "http"),
		mediaServ: mediaServ,
	}
}

func (h *MediaHandler) HandleSaveFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		const op = "MediaHandler.HandleSaveFile"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, _maxFileSize)
		err = c.Request().ParseMultipartForm(_maxFileSize)
		if err != nil {
			logger.Error("failed to parse multipart form", "error", err)

			return echo.ErrBadRequest
		}

		folderName := c.Param("folder")
		fileName := c.Param("file")

		file, fileHeader, err := c.Request().FormFile(fileName)
		if err != nil {
			logger.Error("failed to get file", "error", err)

			return echo.ErrBadRequest
		}
		defer file.Close()

		// TODO: Add content type detection

		obj := storage.Object{
			Data: file,
			Type: fileHeader.Header.Get("Content-Type"),
			Size: fileHeader.Size,
		}

		err = h.mediaServ.SaveFile(ctx, service.SaveFileDTO{
			Folder: folderName,
			File:   fileName,
			Obj:    obj,
		})
		if err != nil {
			logger.Error("failed to save file", "error", err)

			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusCreated)
	}
}
