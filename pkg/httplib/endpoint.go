package httplib

import (
	"errors"
	"net/http"
)

// ? TODO: integrate middleware and endpoint

type Middleware func(next http.Handler) http.Handler

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

func NewMiddlewareFunc(m MiddlewareFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return m(next.ServeHTTP)
	}
}

type Endpoint func(w http.ResponseWriter, r *http.Request) error

type ErroHandler func(w http.ResponseWriter, r *http.Request, err error)

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var (
		code = http.StatusInternalServerError
		data = JSON{"message": "internal server error"}
	)

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		code = apiErr.Code
		data = JSON{"message": apiErr.Message}
	}

	if err := WriteJSON(w, code, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewEndpointWithErroHandler(h Endpoint, errH ErroHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			errH(w, r, err)
		}
	}
}

func NewEndpoint(h Endpoint) http.HandlerFunc {
	return NewEndpointWithErroHandler(h, DefaultErrorHandler)
}
