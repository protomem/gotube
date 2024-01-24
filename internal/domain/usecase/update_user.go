package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserDeps struct {
	Accessor port.UserAccessor
	Mutator  port.UserMutator
}

func UpdateUser(deps UpdateUserDeps) port.UpdateUser {
	return port.UsecaseFunc[port.UpdateUserInput, entity.User](func(
		ctx context.Context,
		input port.UpdateUserInput,
	) (entity.User, error) {
		const op = "usecase.UpdateUser"

		if err := validation.Validate(func(v *validation.Validator) {
			if input.Data.Nickname != nil {
				vrule.Nickname(v, *input.Data.Nickname)
			}
			if input.Data.Email != nil {
				vrule.Email(v, *input.Data.Email)
			}
			if input.Data.NewPassword != nil && input.Data.OldPassword != nil {
				// TODO: different field names
				vrule.Password(v, *input.Data.NewPassword)
				vrule.Password(v, *input.Data.OldPassword)

				if *input.Data.NewPassword == *input.Data.OldPassword {
					v.AddFieldError("newPassword", "must be different from the old password")
				}
			}
			if input.Data.AvatarPath != nil {
				vrule.AvatarPath(v, *input.Data.AvatarPath)
			}
			if (input.Data.NewPassword == nil && input.Data.OldPassword != nil) ||
				(input.Data.NewPassword != nil && input.Data.OldPassword == nil) {
				v.AddError("must provide both old and new password")
			}
			if input.Data.Description != nil {
				vrule.Description(v, *input.Data.Description)
			}
		}); err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		user, err := deps.Accessor.ByNickname(ctx, input.Nickname)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		dto := port.UpdateUserDTO{
			Nickname:    input.Data.Nickname,
			Email:       input.Data.Email,
			AvatarPath:  input.Data.AvatarPath,
			Description: input.Data.Description,
		}

		if input.Data.Email != nil {
			dto.Verified = new(bool)
			*dto.Verified = false
		}

		if input.Data.NewPassword != nil && input.Data.OldPassword != nil {
			if err := bcrypt.CompareHashAndPassword(
				[]byte(user.HashedPassword),
				[]byte(*input.Data.OldPassword),
			); err != nil {
				if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					return entity.User{}, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
				}

				return entity.User{}, fmt.Errorf("%s: %w", op, err)
			}

			hashPass, err := bcrypt.GenerateFromPassword([]byte(*input.Data.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				return entity.User{}, fmt.Errorf("%s: %w", op, err)
			}

			dto.HashedPassword = new(string)
			*dto.HashedPassword = string(hashPass)
		}

		if err := deps.Mutator.Update(ctx, user.ID, dto); err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		newUser, err := deps.Accessor.ByID(ctx, user.ID)
		if err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}

		return newUser, nil
	})
}
