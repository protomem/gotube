package handler

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type User struct {
	logger logging.Logger
	serv   service.User
}

func NewUser(logger logging.Logger, serv service.User) *User {
	return &User{
		logger: logger.With("handler", "user"),
		serv:   serv,
	}
}
