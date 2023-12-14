package main

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/request"
	"github.com/protomem/gotube/internal/response"
	"github.com/protomem/gotube/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
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

	app.mustResponseSend(w, http.StatusCreated, response.Object{"user": user})
}
