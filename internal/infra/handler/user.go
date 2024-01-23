package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/usecase"
	"github.com/protomem/gotube/pkg/request"
	"github.com/protomem/gotube/pkg/response"
)

type User struct {
	*Base

	accessor port.UserAccessor
	mutator  port.UserMutator
}

func NewUser(accessor port.UserAccessor, mutator port.UserMutator) *User {
	return &User{
		Base: NewBase(),

		accessor: accessor,
		mutator:  mutator,
	}
}

func (h *User) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input port.CreateUserInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	deps := usecase.CreateUserDeps{Accessor: h.accessor, Mutator: h.mutator}
	user, err := usecase.CreateUser(deps).Invoke(r.Context(), input)
	if err != nil {
		if entity.IsError(err, entity.ErrUserAlreadyExists) {
			h.ErrorMessage(w, r, http.StatusConflict, entity.ErrUserAlreadyExists.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusCreated, response.Data{"user": user})
}
