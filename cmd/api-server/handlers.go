package main

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/ctxstore"
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

func (app *application) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.UpdateUserByNicknameInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.ByNickname = mustGetUserNicknameFromRequest(r)

	user, err := usecase.UpdateUserByNickname(app.config.baseURL, app.db).Invoke(r.Context(), input)
	if err != nil {
		var vErr *validator.Validator
		if errors.As(err, &vErr) {
			app.failedValidation(w, r, vErr)
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

func (app *application) handleGetVideo(w http.ResponseWriter, r *http.Request) {
	videoID, ok := getVideoIDFromRequest(r)
	if !ok {
		app.badRequest(w, r, errors.New("missing or invalid video ID"))
		return
	}

	requester, isAuth := ctxstore.User(r.Context())

	video, err := usecase.GetVideo(app.db).Invoke(r.Context(), videoID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, model.ErrVideoNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	if video.Public || (isAuth && requester.ID == video.AuthorID) {
		app.mustSendJSON(w, r, http.StatusOK, response.Data{"video": video})
		return
	}

	app.errorMessage(w, r, http.StatusNotFound, model.ErrVideoNotFound.Error(), nil)
}

func (app *application) handleCreateVideo(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateVideoInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	requester := ctxstore.MustUser(r.Context())
	input.AuthorID = requester.ID

	video, err := usecase.CreateVideo(app.config.baseURL, app.db).Invoke(r.Context(), input)
	if err != nil {
		var vErr *validator.Validator
		if errors.As(err, &vErr) {
			app.failedValidation(w, r, vErr)
			return
		}

		if errors.Is(err, model.ErrAlreadyExists) {
			app.errorMessage(w, r, http.StatusConflict, model.ErrVideoAlreadyExists.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustSendJSON(w, r, http.StatusCreated, response.Data{"video": video})
}
