package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/domain/jwt"
	"github.com/protomem/gotube/internal/domain/model"
	"github.com/protomem/gotube/internal/flashstore"
	"github.com/protomem/gotube/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

const (
	_defaultAccessTokenTTL  = 6 * time.Hour
	_defaultRefreshTokenTTL = 3 * 24 * time.Hour

	_tokenIssuer = "gotube"
)

func GetUserByNickname(db *database.DB) Usecase[string, model.User] {
	return UsecaseFunc[string, model.User](func(ctx context.Context, nickname string) (model.User, error) {
		const op = "usecase.GetUserByNickname"

		user, err := db.GetUserByNickname(ctx, nickname)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		return user, nil
	})
}

type (
	CreateUserInput struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
)

func CreateUser(db *database.DB) Usecase[CreateUserInput, model.User] {
	return UsecaseFunc[CreateUserInput, model.User](func(ctx context.Context, input CreateUserInput) (model.User, error) {
		const op = "usecase.CreateUser"

		if err := validator.Validate(func(v *validator.Validator) {
			v.CheckField(validator.MinRunes(input.Nickname, 3), "nickname", "must be at least 3 characters long")
			v.CheckField(validator.MaxRunes(input.Nickname, 20), "nickname", "must be at most 20 characters long")

			v.CheckField(validator.MinRunes(input.Password, 8), "password", "must be at least 8 characters long")
			v.CheckField(validator.MaxRunes(input.Password, 32), "password", "must be at most 32 characters long")

			v.CheckField(validator.IsEmail(input.Email), "email", "must be a valid email")
		}); err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}
		input.Password = string(hashPass)

		id, err := db.InsertUser(ctx, database.InsertUserDTO(input))
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := db.GetUser(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		return user, nil
	})
}

type AuthOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type (
	RegisterInput struct {
		CreateUserInput
	}

	RegisterOutput struct {
		User model.User `json:"user"`
		AuthOutput
	}
)

func Register(authSecret string, db *database.DB, fstore *flashstore.Storage) Usecase[RegisterInput, RegisterOutput] {
	return UsecaseFunc[RegisterInput, RegisterOutput](func(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
		const op = "usecase.Register"

		user, err := CreateUser(db).Invoke(ctx, input.CreateUserInput)
		if err != nil {
			return RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := generateAccessToken(authSecret, _tokenIssuer, user.ID)
		if err != nil {
			return RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			return RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := fstore.PutSession(ctx, model.Session{
			Token:  refreshToken,
			Expiry: time.Now().Add(_defaultRefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return RegisterOutput{
			user,
			AuthOutput{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}, nil
	})
}

type (
	LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginOutput struct {
		User model.User `json:"user"`
		AuthOutput
	}
)

func Login(authSecret string, db *database.DB, fstore *flashstore.Storage) UsecaseFunc[LoginInput, LoginOutput] {
	return UsecaseFunc[LoginInput, LoginOutput](func(ctx context.Context, input LoginInput) (LoginOutput, error) {
		const op = "usecase.Login"

		if err := validator.Validate(func(v *validator.Validator) {
			v.CheckField(validator.IsEmail(input.Email), "email", "must be a valid email")

			v.CheckField(validator.MinRunes(input.Password, 8), "password", "must be at least 8 characters long")
			v.CheckField(validator.MaxRunes(input.Password, 32), "password", "must be at most 32 characters long")
		}); err != nil {
			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := db.GetUserByEmail(ctx, input.Email)
		if err != nil {
			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return LoginOutput{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
			}

			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := generateAccessToken(authSecret, _tokenIssuer, user.ID)
		if err != nil {
			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := fstore.PutSession(ctx, model.Session{
			Token:  refreshToken,
			Expiry: time.Now().Add(_defaultRefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return LoginOutput{
			user,
			AuthOutput{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}, nil
	})
}

type (
	RefreshTokenInput struct {
		RefreshToken string
	}

	RefreshTokenOutput struct {
		AuthOutput
	}
)

func RefreshToken(authSecret string, db *database.DB, fstore *flashstore.Storage) UsecaseFunc[RefreshTokenInput, RefreshTokenOutput] {
	return UsecaseFunc[RefreshTokenInput, RefreshTokenOutput](func(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error) {
		const op = "usecase.RefreshToken"

		session, err := fstore.GetSession(ctx, input.RefreshToken)
		if err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := db.GetUser(ctx, session.UserID)
		if err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := generateAccessToken(authSecret, _tokenIssuer, user.ID)
		if err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := generateRefreshToken()
		if err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := fstore.PutSession(ctx, model.Session{
			Token:  refreshToken,
			Expiry: time.Now().Add(_defaultRefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := fstore.DelSession(ctx, input.RefreshToken); err != nil {
			return RefreshTokenOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return RefreshTokenOutput{
			AuthOutput{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}, nil
	})
}

func Logout(fstore *flashstore.Storage) UsecaseFunc[string, void] {
	return UsecaseFunc[string, void](func(ctx context.Context, token string) (void, error) {
		const op = "usecase.Logout"

		if err := fstore.DelSession(ctx, token); err != nil {
			return void{}, fmt.Errorf("%s: %w", op, err)
		}

		return void{}, nil
	})
}

func generateAccessToken(signingKey, issuer string, userID model.ID) (string, error) {
	accessToken, err := jwt.Generate(jwt.GenerateParams{
		SigningKey: signingKey,
		TTL:        _defaultAccessTokenTTL,
		Subject:    userID,
		Issuer:     issuer,
	})
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func generateRefreshToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
