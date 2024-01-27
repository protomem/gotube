package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/ctxstore"
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

func (h *Comment) getVideoIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return uuid.UUID{}, false
	}

	return videoID, true
}
