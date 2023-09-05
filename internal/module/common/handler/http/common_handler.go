package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
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
		ctx := c.Request().Context()

		res := echo.Map{
			"status": "ok",
		}

		authPayload, ok := jwt.Extract(ctx)
		if ok {
			res["authPayload"] = authPayload
		}

		return c.JSON(http.StatusOK, res)
	}
}
