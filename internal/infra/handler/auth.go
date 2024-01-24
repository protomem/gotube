package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/ctxstore"
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

func (h *Auth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input port.LoginInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	deps := usecase.LoginDeps{
		Conf:     h.conf,
		Accessor: h.accessor,
		Mutator:  h.mutator,
		SessMng:  h.sessMng,
	}
	output, err := usecase.Login(deps).Invoke(r.Context(), input)
	if err != nil {
		var v *validation.Validator
		if errors.As(err, &v) {
			h.FailedValidation(w, r, v)
			return
		}

		if entity.IsError(err, entity.ErrUserNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrUserNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, output)
}

func (h *Auth) HandleRefreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshToken, exists := h.getRefreshTokenFromRequest(r)
	if !exists {
		h.BadRequest(w, r, errors.New("missing refresh token"))
		return
	}

	deps := usecase.RefreshTokensDeps{
		Conf:     h.conf,
		Accessor: h.accessor,
		SessMng:  h.sessMng,
	}
	output, err := usecase.RefreshTokens(deps).
		Invoke(r.Context(), port.RefreshTokensInput{RefreshToken: refreshToken})
	if err != nil {
		// TODO: check other errors
		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, output)
}

func (h *Auth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	refreshToken, exists := h.getRefreshTokenFromRequest(r)
	if !exists {
		h.BadRequest(w, r, errors.New("missing refresh token"))
		return
	}

	deps := usecase.LogoutDeps{SessMng: h.sessMng}
	if _, err := usecase.Logout(deps).
		Invoke(r.Context(), port.LogoutInput{RefreshToken: refreshToken}); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			h.BadRequest(w, r, errors.New("invalid authorization header"))
			return
		}

		accessToken := headerParts[1]

		deps := usecase.VerifyTokenDeps{
			Conf:     h.conf,
			Accessor: h.accessor,
		}
		user, err := usecase.VerifyToken(deps).
			Invoke(r.Context(), port.VerifyTokenInput{AccessToken: accessToken})
		if err != nil {
			if errors.Is(err, port.ErrInvalidToken) {
				h.ErrorMessage(w, r, http.StatusUnauthorized, port.ErrInvalidToken.Error(), nil)
				return
			}

			if entity.IsError(err, entity.ErrUserNotFound) {
				h.ErrorMessage(w, r, http.StatusUnauthorized, entity.ErrUserNotFound.Error(), nil)
				return
			}

			h.ServerError(w, r, err)
			return
		}

		wr := ctxstore.RequestWithUser(r, user)
		next.ServeHTTP(w, wr)
	})
}

func (h *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := ctxstore.User(r.Context())
		if !ok {
			h.ErrorMessage(w, r, http.StatusForbidden, "access denied", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Auth) getRefreshTokenFromRequest(r *http.Request) (string, bool) {
	if r.URL.Query().Has("refresh_token") {
		return r.URL.Query().Get("refresh_token"), true
	} else if r.Header.Get("X-Refresh-Token") != "" {
		return r.Header.Get("X-Refresh-Token"), true
	} else {
		return "", false
	}
}
