package main

import (
	"net/http"

	"github.com/protomem/gotube/internal/response"
)

func (app *application) mustResponseSend(w http.ResponseWriter, status int, data any) {
	if err := response.JSON(w, status, data); err != nil {
		app.serverError(w, nil, err)
		return
	}
}
