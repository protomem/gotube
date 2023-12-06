package middleware

import (
	"net/http"

	"github.com/protomem/gotube/pkg/logging"
)

type MiddlewareFunc func(next http.Handler) http.Handler

type Middlewares struct {
	*Common
}

func New(logger logging.Logger) *Middlewares {
	return &Middlewares{
		Common: NewCommon(logger),
	}
}
