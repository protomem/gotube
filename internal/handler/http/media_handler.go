package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

const _maxObjectSize = 100 * 1024 * 1024 // 100MB

type MediaHandler struct {
	logger logging.Logger
	store  storage.Storage
}

func NewMediaHandler(logger logging.Logger, store storage.Storage) *MediaHandler {
	return &MediaHandler{
		logger: logger.With("handler", "media", "handlerType", "http"),
		store:  store,
	}
}

func (handl *MediaHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.MediaHandler.Get"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		var (
			vars   = mux.Vars(r)
			exists bool
		)

		objParent, exists := vars["parent"]
		if !exists {
			logger.Error("parent missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "parent missing",
			})

			return
		}

		objName, exists := vars["name"]
		if !exists {
			logger.Error("name missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "name missing",
			})

			return
		}

		obj, err := handl.store.Get(ctx, objParent, objName)
		if err != nil {
			logger.Error("failed to get object", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to get object",
			})

			return
		}

		w.Header().Set(httpheader.ContentType, obj.Type)
		w.Header().Set(httpheader.ContentLength, strconv.FormatInt(obj.Size, 10))
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, obj.Src)
	}
}

func (handl *MediaHandler) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.MediaHandler.Get"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		var (
			vars   = mux.Vars(r)
			exists bool
		)

		objParent, exists := vars["parent"]
		if !exists {
			logger.Error("parent missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "parent missing",
			})

			return
		}

		objName, exists := vars["name"]
		if !exists {
			logger.Error("name missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "name missing",
			})

			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, _maxObjectSize)
		err = r.ParseMultipartForm(_maxObjectSize)
		if err != nil {
			logger.Error("failed to parse multipart form", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to parse multipart form",
			})

			return
		}

		file, fileHeader, err := r.FormFile(objName)
		if err != nil {
			logger.Error("failed to get file", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to get file",
			})

			return
		}
		defer func() { _ = file.Close() }()

		err = handl.store.Save(ctx, objParent, objName, storage.Object{
			Size: fileHeader.Size,
			Type: objectTypeResolve(objName),
			Src:  file,
		})
		if err != nil {
			logger.Error("failed to save file", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to save file",
			})

			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (handl *MediaHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.MediaHandler.Delete"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		vars := mux.Vars(r)

		objParent, exists := vars["parent"]
		if !exists {
			logger.Error("parent missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "parent missing",
			})

			return
		}

		objName, exists := vars["name"]
		if !exists {
			logger.Error("name missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "name missing",
			})

			return
		}

		err = handl.store.Delete(ctx, objParent, objName)
		if err != nil {
			logger.Error("failed to delete object", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to delete object",
			})

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func objectTypeResolve(name string) string {
	chunks := strings.Split(name, ".")

	if len(chunks) < 2 {
		return "application/octet-stream"
	}

	switch chunks[len(chunks)-1] {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
