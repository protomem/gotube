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

func (app *application) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input usecase.RegisterInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	output, err := usecase.
		Register(app.config.auth.secret, app.db, app.fstore).
		Invoke(r.Context(), input)
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
