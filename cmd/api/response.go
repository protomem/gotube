package main

import (
	"io"
	"net/http"

	"github.com/protomem/gotube/internal/blobstore"
	"github.com/protomem/gotube/internal/response"
)

func (app *application) mustResponseSend(w http.ResponseWriter, r *http.Request, status int, data any) {
	if err := response.JSON(w, status, data); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) mustSendObject(w http.ResponseWriter, r *http.Request, obj blobstore.Object) {
	w.Header().Set("Content-Type", obj.Type)
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, obj.Body); err != nil {
		app.serverError(w, r, err)
		return
	}
}
