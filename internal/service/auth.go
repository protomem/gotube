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

const (
	_accessTokenTTL  = 6 * time.Hour
	_refreshTokenTTL = 3 * 24 * time.Hour
)

var _ Auth = (*AuthImpl)(nil)

type (
	RegisterDTO struct {
		Nickname string
		Password string
		Email    string
	}

	LoginDTO struct {
		Email    string
		Password string
	}

	Auth interface {
		Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error)
		Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error)
		VerifyToken(ctx context.Context, accessToken string) (model.User, jwt.Payload, error)
		Logout(ctx context.Context, refreshToken string) error
	}

	AuthImpl struct {
		signingKey string
		sessmng    session.Manager
		userServ   User
	}
)

func NewAuth(signingKey string, sessmng session.Manager, userServ User) *AuthImpl {
	return &AuthImpl{
		signingKey: signingKey,
		sessmng:    sessmng,
		userServ:   userServ,
	}
}

func (serv *AuthImpl) Register(ctx context.Context, dto RegisterDTO) (model.User, model.PairTokens, error) {
	const op = "service.Auth.Register"
	var err error

	user, err := serv.userServ.Create(ctx, CreateUserDTO(dto))
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := serv.generatePairTokens(user, _accessTokenTTL)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = serv.sessmng.Set(ctx, session.Session{
		Token:     tokens.Refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(_refreshTokenTTL),
	})
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (serv *AuthImpl) Login(ctx context.Context, dto LoginDTO) (model.User, model.PairTokens, error) {
	const op = "service.Auth.Login"

	user, err := serv.userServ.GetByEmailAndPassword(ctx, dto.Email, dto.Password)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := serv.generatePairTokens(user, _accessTokenTTL)
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = serv.sessmng.Set(ctx, session.Session{
		Token:     tokens.Refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(_refreshTokenTTL),
	})
	if err != nil {
		return model.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (serv *AuthImpl) RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error) {
	const op = "service.Auth.RefreshTokens"
	var err error

	sess, err := serv.sessmng.Get(ctx, refreshToken)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := serv.userServ.Get(ctx, sess.UserID)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := serv.generatePairTokens(user, _accessTokenTTL)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = serv.sessmng.Set(ctx, session.Session{
		Token:     tokens.Refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(_refreshTokenTTL),
	})
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (serv *AuthImpl) VerifyToken(ctx context.Context, accessToken string) (model.User, jwt.Payload, error) {
	const op = "service.Auth.VerifyToken"

	payload, err := serv.parseToken(accessToken)
	if err != nil {
		return model.User{}, payload, fmt.Errorf("%s: %w", op, err)
	}

	user, err := serv.userServ.Get(ctx, payload.UserID)
	if err != nil {
		return model.User{}, payload, fmt.Errorf("%s: %w", op, err)
	}

	return user, payload, nil
}

func (serv *AuthImpl) Logout(ctx context.Context, refreshToken string) error {
	const op = "service.Auth.Logout"

	err := serv.sessmng.Del(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (serv *AuthImpl) generatePairTokens(user model.User, ttl time.Duration) (model.PairTokens, error) {
	const op = "generatePairTokens"
	var err error

	accessToken, err := jwt.Generate(
		jwt.Payload{
			UserID:   user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
			Verified: user.Verified,
		},
		jwt.GenerateParams{
			SigningKey: serv.signingKey,
			TTL:        ttl,
		},
	)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.PairTokens{
		Access:  accessToken,
		Refresh: refreshToken.String(),
	}, nil
}

func (serv *AuthImpl) parseToken(token string) (jwt.Payload, error) {
	const op = "parseToken"

	payload, err := jwt.Parse(token, jwt.ParseParams{
		SigningKey: serv.signingKey,
	})
	if err != nil {
		return jwt.Payload{}, fmt.Errorf("%s: %w", op, err)
	}

	return payload, nil
}
