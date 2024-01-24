package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
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

func (h *User) HandleGet(w http.ResponseWriter, r *http.Request) {
	nickname := h.mustGetNicknameFromRequest(r)

	user, err := h.accessor.ByNickname(r.Context(), nickname)
	if err != nil {
		if entity.IsError(err, entity.ErrUserNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrUserNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"user": user})
}

func (h *User) HandleDelete(w http.ResponseWriter, r *http.Request) {
	nickname := h.mustGetNicknameFromRequest(r)

	user, err := h.accessor.ByNickname(r.Context(), nickname)
	if err != nil {
		if entity.IsError(err, entity.ErrUserNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrUserNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	if err := h.mutator.Delete(r.Context(), user.ID); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *User) mustGetNicknameFromRequest(r *http.Request) string {
	return chi.URLParam(r, "nickname")
}
