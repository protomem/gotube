package main

import (
	"net/http"

	"github.com/protomem/gotube/internal/response"
)

func (app *application) handleStatus(w http.ResponseWriter, r *http.Request) {
	app.mustSendJSON(w, r, http.StatusOK, response.Data{"status": "OK"})
}
