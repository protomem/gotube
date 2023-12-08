package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/protomem/gotube/pkg/header"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type JSON map[string]any

type Header struct {
	Key   string
	Value string
}

func Send(w http.ResponseWriter, code int, res JSON) error {
	return SendWithHeaders(w, code, res)
}

func SendWithHeaders(w http.ResponseWriter, code int, res JSON, headers ...Header) error {
	if res != nil {
		w.Header().Set(header.ContentType, header.ContentTypeJSON)
	}

	for _, h := range headers {
		w.Header().Set(h.Key, h.Value)
	}

	w.WriteHeader(code)

	if res != nil {
		return json.NewEncoder(w).Encode(res)
	} else {
		return nil
	}
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type HandlerFunc func(apiFn APIFunc) http.HandlerFunc

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func BuildHandlerFunc(errFn ErrorHandler) HandlerFunc {
	return func(apiFn APIFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := apiFn(w, r)
			if err != nil {
				errFn(w, r, err)
			}
		}
	}
}

type APIError struct {
	Code int
	Msg  string
	Err  error
}

func NewAPIMsg(code int, msg string) *APIError {
	return &APIError{
		Code: code,
		Msg:  msg,
	}
}

func NewAPIError(code int, msg string, err error) *APIError {
	return &APIError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

func (err *APIError) Error() string {
	if err.Err == nil {
		return err.Msg
	}
	return fmt.Sprintf("%s: %s", err.Msg, err.Err)
}

func (err *APIError) As(target any) bool {
	if _, ok := target.(*APIError); ok {
		return true
	}
	if err.Err != nil {
		return errors.As(err.Err, target)
	}
	return false
}

func DefaultErrorHandler(logger logging.Logger, op string) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var errS error
		defer func() {
			if errS != nil {
				logger.With(
					requestid.Key, requestid.Extract(r.Context()),
					"operation", op,
					"error", errS,
				).Error("failed to send response")
			}
		}()

		var apiErr *APIError
		if errors.As(err, &apiErr) {
			errS = Send(w, apiErr.Code, JSON{"error": apiErr.Error()})
			return
		}

		errS = Send(w, http.StatusInternalServerError, JSON{"error": "internal server error"})
	}
}
