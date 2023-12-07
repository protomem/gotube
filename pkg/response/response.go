package response

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/gotube/pkg/header"
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
	w.Header().Set(header.ContentType, header.ContentTypeJSON)
	for _, h := range headers {
		w.Header().Set(h.Key, h.Value)
	}
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
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
