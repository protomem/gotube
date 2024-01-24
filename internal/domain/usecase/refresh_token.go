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

type RefreshTokensDeps struct {
	Conf     config.Auth
	Accessor port.UserAccessor
	SessMng  port.SessionManager
}

func RefreshTokens(deps RefreshTokensDeps) port.RefreshTokens {
	return port.UsecaseFunc[port.RefreshTokensInput, port.RefreshTokensOutput](func(
		ctx context.Context,
		input port.RefreshTokensInput,
	) (port.RefreshTokensOutput, error) {
		const op = "usecase.RefreshTokens"

		session, err := deps.SessMng.Get(ctx, input.RefreshToken)
		if err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := deps.Accessor.ByID(ctx, session.UserID)
		if err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		accessToken, err := jwt.Generate(jwt.GenerateParams{
			SigningKey: deps.Conf.Secret,
			TTL:        deps.Conf.AccessTokenTTL,
			Subject:    user.ID,
			Issuer:     deps.Conf.Issuer,
		})
		if err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		refreshToken, err := uuid.NewRandom()
		if err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := deps.SessMng.Put(ctx, entity.Session{
			Token:  refreshToken.String(),
			Expiry: time.Now().Add(deps.Conf.RefreshTokenTTL),
			UserID: user.ID,
		}); err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := deps.SessMng.Delete(ctx, session.Token); err != nil {
			return port.RefreshTokensOutput{}, fmt.Errorf("%s: %w", op, err)
		}

		return port.RefreshTokensOutput{
			AccessToken:  accessToken,
			RefreshToken: refreshToken.String(),
		}, nil
	})
}
