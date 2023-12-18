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
	app.mustResponseSend(w, http.StatusOK, response.Object{
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
			app.mustResponseSend(w, http.StatusConflict, response.Object{"error": err.Error()})
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

	app.mustResponseSend(w, http.StatusCreated, response.Object{
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
			app.mustResponseSend(w, http.StatusNotFound, response.Object{"error": err.Error()})
			return
		}

		app.serverError(w, r, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			app.mustResponseSend(w, http.StatusNotFound, response.Object{"error": database.ErrUserNotFound.Error()})
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

	app.mustResponseSend(w, http.StatusOK, response.Object{
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

	app.mustResponseSend(w, http.StatusOK, response.Object{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
