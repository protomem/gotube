package service

import (
	"context"

	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/model"
)

var _ Auth = (*AuthImpl)(nil)

type (
	RegisterDTO struct {
	}

	LoginDTO struct {
	}

	Auth interface {
		Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error)
		Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error)
		VerifyToken(ctx context.Context, accessToken string) (model.User, jwt.Payload, error)
		Logout(ctx context.Context, refreshToken string) error
	}

	AuthImpl struct {
		userServ User
	}
)

func NewAuth(userServ User) *AuthImpl {
	return &AuthImpl{
		userServ: userServ,
	}
}

func (serv *AuthImpl) Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error) {
	panic("unimplemented")
}

func (serv *AuthImpl) Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error) {
	panic("unimplemented")
}

func (serv *AuthImpl) RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error) {
	panic("unimplemented")
}

func (serv *AuthImpl) VerifyToken(ctx context.Context, accessToken string) (model.User, jwt.Payload, error) {
	panic("unimplemented")
}

func (serv *AuthImpl) Logout(ctx context.Context, refreshToken string) error {
	panic("unimplemented")
}
