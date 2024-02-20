package httplib

import (
	"encoding/json"
	"net/http"
)

type JSON map[string]any

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set(HeaderContentType, MIMEApplicationJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}