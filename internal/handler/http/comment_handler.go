package http

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type CommentHandler struct {
	logger logging.Logger
	serv   service.Comment
}

func NewCommentHandler(logger logging.Logger, serv service.Comment) *CommentHandler {
	return &CommentHandler{
		logger: logger.With("handler", "comment", "handlerType", "http"),
		serv:   serv,
	}
}
