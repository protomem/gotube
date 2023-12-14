package handler

import (
	"fmt"
	"net/http"

	"github.com/protomem/gotube/pkg/response"
)

var ErrBadRequest = response.NewAPIMsg(http.StatusBadRequest, "invalid request")

func ErrInternal(msg string, err ...error) *response.APIError {
	err = append(err, nil)
	return response.NewAPIError(http.StatusInternalServerError, msg, err[0])
}

func ErrNotFound(resource string) *response.APIError {
	return response.NewAPIMsg(http.StatusNotFound, fmt.Sprintf("%s not found", resource))
}

func ErrConflict(resource string) *response.APIError {
	return response.NewAPIMsg(http.StatusConflict, fmt.Sprintf("%s already exists", resource))
}
