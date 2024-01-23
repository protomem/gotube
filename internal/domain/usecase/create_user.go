package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserDeps struct {
	Accessor port.UserAccessor
	Mutator  port.UserMutator
}

func CreateUser(deps CreateUserDeps) port.CreateUser {
	return port.UsecaseFunc[port.CreateUserInput, entity.User](func(
		ctx context.Context,
		input port.CreateUserInput,
	) (entity.User, error) {
		const op = "usecase.CreateUser"

		if err := validation.Validate(func(v *validation.Validator) {
			vrule.Nickname(v, input.Nickname)
			vrule.Password(v, input.Password)
			vrule.Email(v, input.Email)
		}); err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		dto := port.InsertUserDTO{
			Nickname: input.Nickname,
			Password: string(hashPass),
			Email:    input.Email,
		}

		id, err := deps.Mutator.Insert(ctx, dto)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := deps.Accessor.ByID(ctx, id)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		return user, nil
	})
}
