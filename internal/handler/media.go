package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/blobstore"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

const _defaultMediaMaxUploadSize = 1024 * 1024 * 100 // 100MB

type Media struct {
	logger logging.Logger
	bstore blobstore.Storage
}

func NewMedia(logger logging.Logger, bstore blobstore.Storage) *Media {
	return &Media{
		logger: logger.With("handler", "media"),
		bstore: bstore,
	}
}

func (h *Media) Get() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		parentName, ok := mux.Vars(r)["parent"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing parent name"})
		}

		fileName, ok := mux.Vars(r)["file"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing file name"})
		}

		obj, err := h.bstore.Get(r.Context(), parentName, fileName)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", obj.Type)

		buff := make([]byte, 1024)
		if _, err := io.CopyBuffer(w, obj.Body, buff); err != nil {
			return err
		}

		return nil
	}, h.errorHandler("handler.Media.Get"))
}

func (h *Media) Save() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		parentName, ok := mux.Vars(r)["parent"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing parent name"})
		}

		fileName, ok := mux.Vars(r)["file"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing file name"})
		}

		r.Body = http.MaxBytesReader(w, r.Body, _defaultMediaMaxUploadSize)
		if err := r.ParseMultipartForm(_defaultMediaMaxUploadSize); err != nil {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "file too large"})
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "could not read file"})
		}
		defer func() { _ = file.Close() }()

		obj := blobstore.Object{
			Type: fileHeader.Header.Get("Content-Type"),
			Size: fileHeader.Size,
			Body: file,
		}

		if err := h.bstore.Put(r.Context(), parentName, fileName, obj); err != nil {
			return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "could not save file"})
		}

		return httplib.SendStatus(w, http.StatusCreated)
	}, h.errorHandler("handler.Media.Save"))
}

func (h *Media) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		parentName, ok := mux.Vars(r)["parent"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing parent name"})
		}

		fileName, ok := mux.Vars(r)["file"]
		if !ok {
			return httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "missing file name"})
		}

		if err := h.bstore.Del(r.Context(), parentName, fileName); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Media.Delete"))
}

func (h *Media) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, blobstore.ErrObjectNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, blobstore.ErrObjectNotFound.Error())
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
