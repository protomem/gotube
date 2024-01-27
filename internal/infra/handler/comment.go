package handler

import (
	"errors"
	"net/http"
	"strconv"

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

type Comment struct {
	*Base

	accessor port.CommentAccessor
	mutator  port.CommentMutator
}

func NewComment(accessor port.CommentAccessor, mutator port.CommentMutator) *Comment {
	return &Comment{
		Base: NewBase(),

		accessor: accessor,
		mutator:  mutator,
	}
}

func (h *Comment) HandleFind(w http.ResponseWriter, r *http.Request) {
	videoID, videoIDOk := h.getVideoIDFromRequest(r)
	if !videoIDOk {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	findOpts, findOptsOk := h.getFindOptionsFromRequest(r)
	if !findOptsOk {
		h.BadRequest(w, r, errors.New("invalid limit or offset"))
		return
	}

	comments, err := h.accessor.AllByVideoID(r.Context(), videoID, findOpts)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Data{"comments": comments})
}

func (h *Comment) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input port.CreateCommentInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}
	input.VideoID = videoID

	requester := ctxstore.MustUser(r.Context())
	input.AuthorID = requester.ID

	deps := usecase.CreateCommentDeps{
		Accessor: h.accessor,
		Mutator:  h.mutator,
	}
	comment, err := usecase.CreateComment(deps).Invoke(r.Context(), input)
	if err != nil {
		var v *validation.Validator
		if errors.As(err, &v) {
			h.FailedValidation(w, r, v)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusCreated, response.Data{"comment": comment})
}

func (h *Comment) HandleDelete(w http.ResponseWriter, r *http.Request) {
	commentID, ok := h.getCommentIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid comment id"))
		return
	}

	if _, err := h.accessor.ByID(r.Context(), commentID); err != nil {
		if entity.IsError(err, entity.ErrCommentNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrCommentNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	if err := h.mutator.Delete(r.Context(), commentID); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Comment) getVideoIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return uuid.UUID{}, false
	}

	return videoID, true
}

func (h *Comment) getCommentIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	commentIDRaw := chi.URLParam(r, "commentId")

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return uuid.UUID{}, false
	}

	return commentID, true
}

func (h *Comment) getFindOptionsFromRequest(r *http.Request) (port.FindOptions, bool) {
	limit, err := strconv.ParseUint(h.DefaultQueryValue(r, "limit", strconv.FormatUint(_defaultLimit, 10)), 10, 64)
	if err != nil {
		return port.FindOptions{}, false
	}

	offset, err := strconv.ParseUint(h.DefaultQueryValue(r, "offset", strconv.FormatUint(_defaultOffset, 10)), 10, 64)
	if err != nil {
		return port.FindOptions{}, false
	}

	return port.FindOptions{Limit: limit, Offset: offset}, true
}
