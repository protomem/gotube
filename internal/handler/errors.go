package handler

import (
	"fmt"
	"net/http"

	"github.com/protomem/gotube/pkg/response"
)

var ErrBadRequest = response.NewAPIMsg(http.StatusBadRequest, "failed to decode request body")

func ErrInternal(msg string) *response.APIError {
	return response.NewAPIMsg(http.StatusInternalServerError, msg)
}

func ErrNotFound(model string) *response.APIError {
	return response.NewAPIMsg(http.StatusNotFound, fmt.Sprintf("%s: not found", model))
}

func ErrConflict(model string) *response.APIError {
	return response.NewAPIMsg(http.StatusConflict, fmt.Sprintf("%s: already exists", model))
}
