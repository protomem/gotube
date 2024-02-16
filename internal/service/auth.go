package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/pkg/jwt"
)

const _defaultTokenIssuer = "gotube"

var _ Auth = (*AuthImpl)(nil)

type (
	LoginDTO struct {
		Email    string
		Password string
	}
)

type (
	Auth interface {
		Login(ctx context.Context, dto LoginDTO) (token string, user model.User, err error)
		Verify(ctx context.Context, token string) (model.User, error)
	}

	AuthImpl struct {
		conf     config.Auth
		userServ User
	}
)

func NewAuth(conf config.Auth, userServ User) *AuthImpl {
	return &AuthImpl{
		conf:     conf,
		userServ: userServ,
	}
}

func (s *AuthImpl) Login(ctx context.Context, dto LoginDTO) (string, model.User, error) {
	const op = "service.Auth.Login"

	user, err := s.userServ.GetByEmailAndPassword(ctx, dto.Email, dto.Password)
	if err != nil {
		return "", model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.Generate(jwt.GenerateParams{
		SigningKey: s.conf.Secret,
		TTL:        s.conf.AccessTokenTTL,
		Subject:    user.Nickname,
		Issuer:     _defaultTokenIssuer,
	})
	if err != nil {
		return "", model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return token, user, nil
}

func (s *AuthImpl) Verify(ctx context.Context, token string) (model.User, error) {
	const op = "service.Auth.Verify"

	subject, err := jwt.Parse(token, jwt.ParseParams{
		SigningKey: s.conf.Secret,
		Issuer:     _defaultTokenIssuer,
	})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userServ.GetByNickname(ctx, subject)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
