package service

import "github.com/protomem/gotube/internal/repository"

var _ User = (*UserImpl)(nil)

type (
	User interface{}

	UserImpl struct {
		repo repository.User
	}
)

func NewUser(repo repository.User) *UserImpl {
	return &UserImpl{
		repo: repo,
	}
}
