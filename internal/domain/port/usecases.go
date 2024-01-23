package port

import "github.com/protomem/gotube/internal/domain/entity"

type (
	CreateUserInput struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
)

type CreateUser = Usecase[CreateUserInput, entity.User]

type (
	RegisterInput = CreateUserInput

	RegisterOutput struct {
		User         entity.User `json:"user"`
		RefreshToken string      `json:"refreshToken"`
		AccessToken  string      `json:"accessToken"`
	}
)

type Register = Usecase[RegisterInput, RegisterOutput]

type (
	LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginOutput struct {
		User         entity.User `json:"user"`
		RefreshToken string      `json:"refreshToken"`
		AccessToken  string      `json:"accessToken"`
	}
)

type Login = Usecase[LoginInput, LoginOutput]
