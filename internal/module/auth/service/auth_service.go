package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/auth/model"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
	userserv "github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/storage"
)

const (
	AccessTokenTTL  = 6 * time.Hour
	RefreshTokenTTL = 30 * 24 * time.Hour
)

var _ AuthService = (*AuthServiceImpl)(nil)

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
	AuthService interface {
		Register(ctx context.Context, dto RegisterDTO) (usermodel.User, model.PairTokens, error)
		Login(ctx context.Context, dto LoginDTO) (usermodel.User, model.PairTokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error)
		VerifyToken(ctx context.Context, accessToken string) (usermodel.User, jwt.Payload, error)
	}

	AuthServiceImpl struct {
		secret string

		sessmng storage.SessionManager

		userServ userserv.UserService
	}
)

func NewAuthService(secret string, sessmng storage.SessionManager, userServ userserv.UserService) *AuthServiceImpl {
	return &AuthServiceImpl{
		secret:   secret,
		sessmng:  sessmng,
		userServ: userServ,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, dto RegisterDTO) (usermodel.User, model.PairTokens, error) {
	const op = "AuthService.Register"
	var err error

	user, err := s.userServ.CreateUser(ctx, userserv.CreateUserDTO(dto))
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.genTokens(user)
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.sessmng.SetSession(ctx, tokens.RefreshToken, storage.Session{
		UserID:    user.ID.String(),
		ExpiredAt: time.Now().Add(RefreshTokenTTL),
	})
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, dto LoginDTO) (usermodel.User, model.PairTokens, error) {
	const op = "AuthService.Login"
	var err error

	user, err := s.userServ.FindOneUserByEmailAndPassword(ctx, userserv.FindOneUserByEmailAndPasswordDTO(dto))
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.genTokens(user)
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.sessmng.SetSession(ctx, tokens.RefreshToken, storage.Session{
		UserID:    user.ID.String(),
		ExpiredAt: time.Now().Add(RefreshTokenTTL),
	})
	if err != nil {
		return usermodel.User{}, model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, tokens, nil
}

func (s *AuthServiceImpl) RefreshTokens(ctx context.Context, refreshToken string) (model.PairTokens, error) {
	const op = "AuthService.RefreshTokens"
	var err error

	sess, err := s.sessmng.GetSession(ctx, refreshToken)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := uuid.Parse(sess.UserID)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userServ.FindOneUser(ctx, userID)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := s.genTokens(user)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.sessmng.SetSession(ctx, tokens.RefreshToken, storage.Session{
		UserID:    user.ID.String(),
		ExpiredAt: time.Now().Add(RefreshTokenTTL),
	})
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.sessmng.DelSession(ctx, refreshToken)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (s *AuthServiceImpl) VerifyToken(
	ctx context.Context,
	accessToken string,
) (usermodel.User, jwt.Payload, error) {
	const op = "AuthService.VerifyTokens"
	var err error

	params := jwt.ParseParams{SigningKey: s.secret}
	payload, err := jwt.Parse(accessToken, params)
	if err != nil {
		return usermodel.User{}, jwt.Payload{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userServ.FindOneUser(ctx, payload.UserID)
	if err != nil {
		return usermodel.User{}, jwt.Payload{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, payload, nil
}

func (s *AuthServiceImpl) genTokens(user usermodel.User) (model.PairTokens, error) {
	const op = "generate tokens"
	var err error

	payload := jwt.Payload{UserID: user.ID, Nickname: user.Nickname, Email: user.Email, Verified: user.Verified}
	params := jwt.GenerateParams{SigningKey: s.secret, TTL: AccessTokenTTL}
	accessToken, err := jwt.Generate(payload, params)
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return model.PairTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.PairTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}
