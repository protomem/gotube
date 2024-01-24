package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/port"
)

type LogoutDeps struct {
	SessMng port.SessionManager
}

func Logout(deps LogoutDeps) port.Logout {
	return port.UsecaseFunc[port.LogoutInput, port.Void](func(
		ctx context.Context,
		input port.LogoutInput,
	) (port.Void, error) {
		const op = "usecase.Logout"

		if err := deps.SessMng.Delete(ctx, input.RefreshToken); err != nil {
			return port.Void{}, fmt.Errorf("%s: %w", op, err)
		}

		return port.Void{}, nil
	})
}
