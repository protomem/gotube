package main

import (
	"errors"
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
	w.Header().Set(HeaderContentType, obj.Type)
	w.Header().Set(HeaderXContentType, "nosniff")
	w.Header().Set(HeaderConnection, "keep-alive")
	w.Header().Set(HeaderTransferEncoding, "chunked")

	wc := http.NewResponseController(w)

	sending := true
	for sending {
		if _, err := io.Copy(w, obj.Body); err != nil {
			if errors.Is(err, io.EOF) {
				sending = false
			} else {
				app.serverError(w, r, err)
				return
			}
		}

		if err := wc.Flush(); err != nil {
			app.serverError(w, r, err)
			return
		}
	}
}
