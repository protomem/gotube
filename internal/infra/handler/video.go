package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/usecase"
	"github.com/protomem/gotube/pkg/request"
	"github.com/protomem/gotube/pkg/response"
	"github.com/protomem/gotube/pkg/validation"
)

type Video struct {
	*Base

	accessor port.VideoAccessor
	mutator  port.VideoMutator
}

func NewVideo(accessor port.VideoAccessor, mutator port.VideoMutator) *Video {
	return &Video{
		Base: NewBase(),

		accessor: accessor,
		mutator:  mutator,
	}
}

func (h *Video) HandleGet(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	video, err := h.accessor.ByID(r.Context(), videoID)
	if err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	requester, isAuth := ctxstore.User(r.Context())
	if !video.Public && (!isAuth || requester.ID != video.Author.ID) {
		h.ErrorMessage(w, r, http.StatusForbidden, "access denied", nil)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"video": video})
}

func (h *Video) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input port.CreateVideoInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	requester := ctxstore.MustUser(r.Context())
	input.AuthorID = requester.ID

	deps := usecase.CreateVideoDeps{
		Accessor: h.accessor,
		Mutator:  h.mutator,
	}
	video, err := usecase.CreateVideo(deps).Invoke(r.Context(), input)
	if err != nil {
		var v *validation.Validator
		if errors.As(err, &v) {
			h.FailedValidation(w, r, v)
			return
		}

		if entity.IsError(err, entity.ErrVideoAlreadyExists) {
			h.ErrorMessage(w, r, http.StatusConflict, entity.ErrVideoAlreadyExists.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusCreated, response.Data{"video": video})
}

func (h *Video) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	var input port.UpdateVideoData
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	deps := usecase.UpdateVideoDeps{
		Accessor: h.accessor,
		Mutator:  h.mutator,
	}
	video, err := usecase.UpdateVideo(deps).Invoke(r.Context(), port.UpdateVideoInput{
		ID:   videoID,
		Data: input,
	})
	if err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"video": video})
}

func (h *Video) HandleDelete(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	if _, err := h.accessor.ByID(r.Context(), videoID); err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	if err := h.mutator.Delete(r.Context(), videoID); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Video) getVideoIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return uuid.UUID{}, false
	}

	return videoID, true
}
