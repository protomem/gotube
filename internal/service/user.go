package service

import (
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/hashing"
)

var _ User = (*UserImpl)(nil)

type (
	User interface{}

	UserImpl struct {
		repo   repository.User
		hasher hashing.Hasher
	}
)

func NewUser(repo repository.User, hasher hashing.Hasher) *UserImpl {
	return &UserImpl{
		repo:   repo,
		hasher: hasher,
	}
}
