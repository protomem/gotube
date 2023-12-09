package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/session"
)

const _accessTokenTTL = 6 * time.Hour

var _ Auth = (*AuthImpl)(nil)

type RegisterDTO struct {
	Nickname string
	Password string
	Email    string
}

type LoginDTO struct {
	Email    string
	Password string
}

type (
	Auth interface {
		Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error)
		Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error)
	}

	AuthImpl struct {
		authSecret string
		userServ   User
		sessmng    session.Manager
	}
)

func NewAuth(authSecret string, userServ User, sessmng session.Manager) *AuthImpl {
	return &AuthImpl{
		authSecret: authSecret,
		userServ:   userServ,
		sessmng:    sessmng,
	}
}

func (s *AuthImpl) Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error) {
	const op = "service:Auth.Register"

	user, err := s.userServ.Create(ctx, CreateUserDTO(dto))
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (s *AuthImpl) Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error) {
	const op = "service:Auth.Login"

	user, err := s.userServ.GetByEmailAndPassword(ctx, dto.Email, dto.Password)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (s *AuthImpl) generateTokens(user model.User) (model.PairTokens, error) {
	const op = "generate tokens"

	payload := jwt.Payload{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
		Verified: user.Verified,
	}

	params := jwt.GenerateParams{SigningKey: s.authSecret, TTL: _accessTokenTTL}

	accessToken, err := jwt.Generate(payload, params)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.PairTokens{Access: accessToken, Refresh: refreshToken.String()}, nil
}
