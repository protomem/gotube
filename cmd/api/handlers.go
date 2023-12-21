package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/cookies"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/flashstore"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/request"
	"github.com/protomem/gotube/internal/response"
	"github.com/protomem/gotube/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrMissingRefreshToken = errors.New("missing refresh token")

func (app *application) handleStatus(w http.ResponseWriter, r *http.Request) {
	app.mustResponseSend(w, r, http.StatusOK, response.Object{
		"status": "OK",
	})
}

func (app *application) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`

		validator.Validator `json:"-"`
	}

	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(
		validator.MinRunes(input.Nickname, 3) && validator.MaxRunes(input.Nickname, 18),
		"nickname", "nickname must be between 3 and 18 characters",
	)
	input.Validator.CheckField(
		validator.MinRunes(input.Password, 8) && validator.MaxRunes(input.Password, 16),
		"password", "password must be between 8 and 16 characters",
	)
	input.Validator.CheckField(
		validator.NotBlank(input.Email) &&
			validator.IsEmail(input.Email) && validator.MaxRunes(input.Email, 32),
		"email", "invalid email address",
	)

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	dto := database.InsertUserDTO{
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userID, err := app.db.InsertUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			app.errorMessage(w, r, http.StatusConflict, database.ErrUserAlreadyExists.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	user, err := app.db.GetUser(r.Context(), userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	accessToken, err := jwt.Generate(jwt.GenerateParams{
		SigningKey: app.config.auth.secretKey,
		TTL:        app.config.auth.tokenTTL,
		Subject:    user.ID,
		Issuer:     app.config.baseURL,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.fstore.PutSession(r.Context(), flashstore.Session{
		Token:  refreshToken.String(),
		TTL:    app.config.auth.sessionTTL,
		UserID: user.ID,
	}); err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := cookies.WriteSigned(w, http.Cookie{
		Name:     "session",
		Value:    refreshToken.String(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(app.config.auth.sessionTTL), // TODO: to seconds
	}, app.config.cookie.secretKey); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusCreated, response.Object{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`

		validator.Validator `json:"-"`
	}

	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(
		validator.NotBlank(input.Email) && validator.IsEmail(input.Email),
		"email", "invalid email address",
	)
	input.Validator.CheckField(
		validator.MinRunes(input.Password, 8) && validator.MaxRunes(input.Password, 16),
		"password", "password must be between 8 and 16 characters",
	)

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	user, err := app.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	accessToken, err := jwt.Generate(jwt.GenerateParams{
		SigningKey: app.config.auth.secretKey,
		TTL:        app.config.auth.tokenTTL,
		Subject:    user.ID,
		Issuer:     app.config.baseURL,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.fstore.PutSession(r.Context(), flashstore.Session{
		Token:  refreshToken.String(),
		TTL:    app.config.auth.sessionTTL,
		UserID: user.ID,
	}); err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := cookies.WriteSigned(w, http.Cookie{
		Name:     "session",
		Value:    refreshToken.String(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(app.config.auth.sessionTTL), // TODO: to seconds
	}, app.config.cookie.secretKey); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusOK, response.Object{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func (app *application) handleLogout(w http.ResponseWriter, r *http.Request) {
	refreshToken := getRefreshTokenFromRequest(r, app.config.cookie.secretKey)
	if refreshToken == "" {
		app.serverError(w, r, ErrMissingRefreshToken)
		return
	}

	if err := app.fstore.DelSession(r.Context(), refreshToken); err != nil {
		app.serverError(w, r, err)
		return
	}

	// TODO: remove session cookie

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenFromRequest := getRefreshTokenFromRequest(r, app.config.cookie.secretKey)
	if refreshTokenFromRequest == "" {
		app.serverError(w, r, ErrMissingRefreshToken)
		return
	}

	session, err := app.fstore.GetSession(r.Context(), refreshTokenFromRequest)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.db.GetUser(r.Context(), session.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	accessToken, err := jwt.Generate(jwt.GenerateParams{
		SigningKey: app.config.auth.secretKey,
		TTL:        app.config.auth.tokenTTL,
		Subject:    user.ID,
		Issuer:     app.config.baseURL,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.fstore.PutSession(r.Context(), flashstore.Session{
		Token:  refreshToken.String(),
		TTL:    app.config.auth.sessionTTL,
		UserID: user.ID,
	}); err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.fstore.DelSession(r.Context(), refreshTokenFromRequest); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusOK, response.Object{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (app *application) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userNickname := getUserNicknameFromRequest(r)

	user, err := app.db.GetUserByNickname(r.Context(), userNickname)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusOK, response.Object{"user": user})
}

func (app *application) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname    *string `json:"nickname"`
		Email       *string `json:"email"`
		AvatarPath  *string `json:"avatarPath"`
		Description *string `json:"description"`

		OldPassword *string `json:"oldPassword"`
		NewPassword *string `json:"newPassword"`

		validator.Validator `json:"-"`
	}

	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.serverError(w, r, err)
		return
	}

	if input.Nickname != nil {
		input.Validator.CheckField(
			validator.MinRunes(*input.Nickname, 3) && validator.MaxRunes(*input.Nickname, 18),
			"nickname", "nickname must be between 3 and 18 characters",
		)
	}
	if input.Email != nil {
		input.Validator.CheckField(
			validator.NotBlank(*input.Email) &&
				validator.IsEmail(*input.Email) && validator.MaxRunes(*input.Email, 32),
			"email", "invalid email address",
		)
	}
	if input.AvatarPath != nil {
		input.Validator.CheckField(
			validator.MaxRunes(*input.AvatarPath, 255) && validator.IsURL(*input.AvatarPath),
			"avatarPath", "invalid avatar path",
		)
	}
	if input.Description != nil {
		input.Validator.CheckField(
			validator.MaxRunes(*input.Description, 500),
			"description", "description must be less than 500 characters",
		)
	}
	if input.NewPassword != nil && input.OldPassword != nil {
		input.Validator.CheckField(
			validator.MinRunes(*input.NewPassword, 8) && validator.MaxRunes(*input.NewPassword, 32),
			"newPassword", "new password must be between 8 and 32 characters",
		)
		input.Validator.CheckField(
			validator.MinRunes(*input.OldPassword, 8) && validator.MaxRunes(*input.OldPassword, 32),
			"oldPassword", "old password must be between 8 and 32 characters",
		)
		input.Validator.Check(*input.NewPassword != *input.OldPassword, "new password must be different from old password")
	}

	if input.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	userNickname := getUserNicknameFromRequest(r)

	oldUser, err := app.db.GetUserByNickname(r.Context(), userNickname)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	var hashedNewPassword *string
	if input.NewPassword != nil && input.OldPassword != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(*input.OldPassword)); err != nil {
			app.serverError(w, r, errors.New("invalid old password"))
			return
		}

		hashedNewPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(*input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		hashedNewPassword = new(string)
		*hashedNewPassword = string(hashedNewPasswordBytes)
	}

	dto := database.UpdateUserDTO{
		Nickname:    input.Nickname,
		Password:    hashedNewPassword,
		Email:       input.Email,
		AvatarPath:  input.AvatarPath,
		Description: input.Description,
	}

	if err := app.db.UpdateUser(r.Context(), oldUser.ID, dto); err != nil {
		app.serverError(w, r, err)
		return
	}

	newUser, err := app.db.GetUser(r.Context(), oldUser.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusOK, response.Object{"user": newUser})
}

func (app *application) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userNickname := getUserNicknameFromRequest(r)

	user, err := app.db.GetUserByNickname(r.Context(), userNickname)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrUserNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	if err := app.db.DeleteUser(r.Context(), user.ID); err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleGetVideo(w http.ResponseWriter, r *http.Request) {
	videoID, err := getVideoIDFromRequest(r)
	if err != nil {
		app.badRequest(w, r, errors.New("invalid video id"))
		return
	}

	video, err := app.db.GetVideo(r.Context(), videoID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.errorMessage(w, r, http.StatusNotFound, database.ErrVideoNotFound.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	user, isAuth := contextGetUser(r)
	if !video.Public && (!isAuth && video.AuthorID != user.ID) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	app.mustResponseSend(w, r, http.StatusOK, response.Object{"video": video})
}

func (app *application) handleCreateVideo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		ThumbnailPath *string `json:"thumbnailPath"`
		VideoPath     *string `json:"videoPath"`
		Public        *bool   `json:"isPublic"`

		validator.Validator `json:"-"`
	}

	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(
		validator.NotBlank(input.Title) &&
			validator.MinRunes(input.Title, 3) && validator.MaxRunes(input.Title, 100),
		"title", "title must be between 3 and 100 characters",
	)
	if input.Description != nil {
		input.Validator.CheckField(
			validator.NotBlank(*input.Description) && validator.MaxRunes(*input.Description, 1000),
			"description", "description must be between 1 and 1000 characters",
		)
	}
	if input.ThumbnailPath != nil {
		input.Validator.CheckField(
			validator.NotBlank(*input.ThumbnailPath) && validator.IsURL(*input.ThumbnailPath),
			"thumbnailPath", "invalid thumbnail path",
		)
	}
	if input.VideoPath != nil {
		input.Validator.CheckField(
			validator.NotBlank(*input.VideoPath) && validator.IsURL(*input.VideoPath),
			"videoPath", "invalid video path",
		)
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	user, _ := contextGetUser(r)

	dto := database.InsertVideoDTO{
		Title:    input.Title,
		AuthorID: user.ID,
	}
	if input.Description != nil {
		dto.Description = *input.Description
	}
	if input.ThumbnailPath != nil {
		dto.ThumbnailPath = *input.ThumbnailPath
	}
	if input.VideoPath != nil {
		dto.VideoPath = *input.VideoPath
	}
	if input.Public != nil {
		dto.Public = *input.Public
	}

	videoID, err := app.db.InsertVideo(r.Context(), dto)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			app.errorMessage(w, r, http.StatusConflict, database.ErrVideoAlreadyExists.Error(), nil)
			return
		}

		app.serverError(w, r, err)
		return
	}

	video, err := app.db.GetVideo(r.Context(), videoID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.mustResponseSend(w, r, http.StatusCreated, response.Object{"video": video})
}
