package handler

import (
	"errors"
	"net/http"

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
