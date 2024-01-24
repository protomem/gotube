package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/jwt"
)

type VerifyTokenDeps struct {
	Conf     config.Auth
	Accessor port.UserAccessor
}

func VerifyToken(deps VerifyTokenDeps) port.VerifyToken {
	return port.UsecaseFunc[port.VerifyTokenInput, entity.User](func(
		ctx context.Context,
		input port.VerifyTokenInput,
	) (entity.User, error) {
		const op = "usecase.VerifyToken"

		userID, err := jwt.Parse(input.AccessToken, jwt.ParseParams{
			SigningKey: deps.Conf.Secret,
			Issuer:     deps.Conf.Issuer,
		})
		if err != nil {
			if errors.Is(err, jwt.ErrInvalidToken) {
				return entity.User{}, fmt.Errorf("%s: %w", op, port.ErrInvalidToken)
			}

			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := deps.Accessor.ByID(ctx, userID)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		return user, nil
	})
}
