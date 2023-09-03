package service

import "github.com/protomem/gotube/internal/module/user/repository"

var _ UserService = (*UserServiceImpl)(nil)

type (
	UserService interface{}

	UserServiceImpl struct {
		userRepo repository.UserRepository
	}
)

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}
