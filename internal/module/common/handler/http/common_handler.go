package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/pkg/logging"
)

type CommonHandler struct {
	logger logging.Logger
}

func NewCommonHandler(logger logging.Logger) *CommonHandler {
	return &CommonHandler{
		logger: logger.With("handler", "common", "handlerType", "http"),
	}
}

func (h *CommonHandler) HandleHealthCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"status": "ok",
		})
	}
}
