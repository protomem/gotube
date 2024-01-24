package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginDeps struct {
	Conf     config.Auth
	Accessor port.UserAccessor
	Mutator  port.UserMutator
	SessMng  port.SessionManager
}

func Login(deps LoginDeps) port.Login {
	return port.UsecaseFunc[port.LoginInput, port.LoginOutput](func(
		ctx context.Context,
		input port.LoginInput,
	) (port.LoginOutput, error) {
		const op = "usecase.Login"

		if err := validation.Validate(func(v *validation.Validator) {
			vrule.Email(v, input.Email)
			vrule.Password(v, input.Password)
		}); err != nil {
			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := deps.Accessor.ByEmail(ctx, input.Email)
		if err != nil {
			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return port.LoginOutput{}, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
			}

			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := jwt.Generate(jwt.GenerateParams{
			SigningKey: deps.Conf.Secret,
			TTL:        deps.Conf.AccessTokenTTL,
			Subject:    user.ID,
			Issuer:     deps.Conf.Issuer,
		})
		if err != nil {
			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := uuid.NewRandom()
		if err != nil {
			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := deps.SessMng.Put(ctx, entity.Session{
			Token:  refreshToken.String(),
			Expiry: time.Now().Add(deps.Conf.RefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return port.LoginOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return port.LoginOutput{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken.String(),
		}, nil
	})
}
