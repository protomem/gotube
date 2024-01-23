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
