package main

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/domain/model"
	"github.com/protomem/gotube/internal/domain/usecase"
	"github.com/protomem/gotube/internal/request"
	"github.com/protomem/gotube/internal/response"
	"github.com/protomem/gotube/internal/validator"
)

func (app *application) handleStatus(w http.ResponseWriter, r *http.Request) {
	app.mustSendJSON(w, r, http.StatusOK, response.Data{"status": "OK"})
}

func (app *application) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userNickname := mustGetUserNicknameFromRequest(r)

	user, err := usecase.GetUserByNickname(app.db).Invoke(r.Context(), userNickname)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, model.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustSendJSON(w, r, http.StatusOK, response.Data{"user": user})
}

func (app *application) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userNickname := mustGetUserNicknameFromRequest(r)

	if _, err := usecase.DeleteUserByNickname(app.db).Invoke(r.Context(), userNickname); err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input usecase.RegisterInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	output, err := usecase.Register(app.config.auth.secret, app.db, app.fstore).Invoke(r.Context(), input)
	if err != nil {
		var vErr *validator.Validator
		if errors.As(err, &vErr) {
			app.failedValidation(w, r, vErr)
			return
		}

		if errors.Is(err, model.ErrAlreadyExists) {
			app.errorMessage(w, r, http.StatusConflict, model.ErrUserAlreadyExists.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustSendJSON(w, r, http.StatusCreated, output)
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input usecase.LoginInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	output, err := usecase.Login(app.config.auth.secret, app.db, app.fstore).Invoke(r.Context(), input)
	if err != nil {
		var vErr *validator.Validator
		if errors.As(err, &vErr) {
			app.failedValidation(w, r, vErr)
			return
		}

		if errors.Is(err, model.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, model.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustSendJSON(w, r, http.StatusOK, output)
}

func (app *application) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := getRefreshTokenFromRequest(r)
	if !ok {
		app.badRequest(w, r, errors.New("missing refresh token"))
		return
	}

	output, err := usecase.
		RefreshToken(app.config.auth.secret, app.db, app.fstore).
		Invoke(r.Context(), usecase.RefreshTokenInput{RefreshToken: refreshToken})
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, model.ErrSessionNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustSendJSON(w, r, http.StatusOK, output)
}

func (app *application) handleLogout(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := getRefreshTokenFromRequest(r)
	if !ok {
		app.badRequest(w, r, errors.New("missing refresh token"))
		return
	}

	if _, err := usecase.Logout(app.fstore).Invoke(r.Context(), refreshToken); err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
