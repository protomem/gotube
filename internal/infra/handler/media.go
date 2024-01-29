package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/infra/blobstore"
)

const _defaultMediaMaxUploadSize = 1024 * 1024 * 100

type Media struct {
	*Base

	bstore *blobstore.Storage
}

func NewMedia(bstore *blobstore.Storage) *Media {
	return &Media{
		Base: NewBase(),

		bstore: bstore,
	}
}

func (h *Media) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	parentName := h.baseParentName(h.mustGetParentNameFromRequest(r))
	fileName := h.mustGetFileNameFromRequest(r)

	obj, err := h.bstore.GetObject(r.Context(), parentName, fileName)
	if err != nil {
		if errors.Is(err, blobstore.ErrObjectNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, blobstore.ErrObjectNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", obj.Type)

	buff := make([]byte, 1024)
	if _, err := io.CopyBuffer(w, obj.Body, buff); err != nil {
		h.ServerError(w, r, err)
		return
	}
}

func (h *Media) HandleSaveFile(w http.ResponseWriter, r *http.Request) {
	parentName := h.baseParentName(h.mustGetParentNameFromRequest(r))
	fileName := h.mustGetFileNameFromRequest(r)

	r.Body = http.MaxBytesReader(w, r.Body, _defaultMediaMaxUploadSize)
	if err := r.ParseMultipartForm(_defaultMediaMaxUploadSize); err != nil {
		h.BadRequest(w, r, errors.New("file too large"))
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.BadRequest(w, r, err)
		return
	}
	defer func() { _ = file.Close() }()

	obj := blobstore.Object{
		Type: fileHeader.Header.Get("Content-Type"),
		Size: fileHeader.Size,
		Body: file,
	}

	if err := h.bstore.PutObject(r.Context(), parentName, fileName, obj); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Media) HandleRemoveFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Media) mustGetParentNameFromRequest(r *http.Request) string {
	return chi.URLParam(r, "parentName")
}

func (h *Media) mustGetFileNameFromRequest(r *http.Request) string {
	return chi.URLParam(r, "fileName")
}

func (h *Media) baseParentName(parentName string) string {
	return fmt.Sprintf("media-%s", parentName)
}
