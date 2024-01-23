package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/jwt"
)

type RegisterDeps struct {
	Conf     config.Auth
	Accessor port.UserAccessor
	Mutator  port.UserMutator
	SessMng  port.SessionManager
}

func Register(deps RegisterDeps) port.Register {
	return port.UsecaseFunc[port.RegisterInput, port.RegisterOutput](func(
		ctx context.Context,
		input port.RegisterInput,
	) (port.RegisterOutput, error) {
		const op = "usecase.Register"

		user, err := CreateUser(CreateUserDeps{
			Accessor: deps.Accessor,
			Mutator:  deps.Mutator,
		}).Invoke(ctx, input)
		if err != nil {
			return port.RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := jwt.Generate(jwt.GenerateParams{
			SigningKey: deps.Conf.Secret,
			TTL:        deps.Conf.AccessTokenTTL,
			Subject:    user.ID,
			Issuer:     user.Nickname,
		})
		if err != nil {
			return port.RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := uuid.NewRandom()
		if err != nil {
			return port.RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := deps.SessMng.Put(ctx, entity.Session{
			Token:  refreshToken.String(),
			Expiry: time.Now().Add(deps.Conf.RefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return port.RegisterOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return port.RegisterOutput{
			User:         user,
			RefreshToken: refreshToken.String(),
			AccessToken:  accessToken,
		}, nil
	})
}
