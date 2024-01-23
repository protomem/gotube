package handler

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/usecase"
	"github.com/protomem/gotube/pkg/request"
	"github.com/protomem/gotube/pkg/validation"
)

type Auth struct {
	*Base

	conf     config.Auth
	accessor port.UserAccessor
	mutator  port.UserMutator
	sessMng  port.SessionManager
}

func NewAuth(
	conf config.Config,
	accessor port.UserAccessor,
	mutator port.UserMutator,
	sessMng port.SessionManager,
) *Auth {
	return &Auth{
		Base: NewBase(),

		conf:     conf.Auth,
		accessor: accessor,
		mutator:  mutator,
		sessMng:  sessMng,
	}
}

func (h *Auth) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var input port.RegisterInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	deps := usecase.RegisterDeps{
		Conf:     h.conf,
		Accessor: h.accessor,
		Mutator:  h.mutator,
		SessMng:  h.sessMng,
	}
	output, err := usecase.Register(deps).Invoke(r.Context(), input)
	if err != nil {
		var v *validation.Validator
		if errors.As(err, &v) {
			h.FailedValidation(w, r, v)
			return
		}

		if entity.IsError(err, entity.ErrUserAlreadyExists) {
			h.ErrorMessage(w, r, http.StatusConflict, entity.ErrUserAlreadyExists.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusCreated, output)
}
