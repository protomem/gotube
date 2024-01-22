package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/pkg/response"
	"github.com/protomem/gotube/pkg/validation"
)

type MiddlewareFunc func(next http.Handler) http.Handler

type Base struct{}

func NewBase() *Base {
	return &Base{}
}

func (h *Base) ReportError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url, "trace", trace)
	if logger, ok := ctxstore.Logger(r.Context()); ok {
		logger.Error(message, requestAttrs)
	}
}

func (h *Base) ErrorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	if err := response.JSONWithHeaders(w, status, response.Data{"error": message}, headers); err != nil {
		h.ReportError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Base) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.ReportError(r, err)

	message := "The server encountered a problem and could not process your request"
	h.ErrorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (h *Base) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	h.ErrorMessage(w, r, http.StatusNotFound, message, nil)
}

func (h *Base) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	h.ErrorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (h *Base) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (h *Base) FailedValidation(w http.ResponseWriter, r *http.Request, v *validation.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		h.ServerError(w, r, err)
	}
}

func (h *Base) MustSendJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	if err := response.JSON(w, status, data); err != nil {
		h.ServerError(w, r, err)
	}
}
