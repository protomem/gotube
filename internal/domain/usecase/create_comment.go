package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/pkg/validation"
)

type CreateCommentDeps struct {
	Accessor port.CommentAccessor
	Mutator  port.CommentMutator
}

func CreateComment(deps CreateCommentDeps) port.CreateComment {
	return port.UsecaseFunc[port.CreateCommentInput, entity.Comment](func(
		ctx context.Context,
		input port.CreateCommentInput,
	) (entity.Comment, error) {
		const op = "usecase.CreateComment"

		if err := validation.Validate(func(v *validation.Validator) {
			vrule.Content(v, input.Content)
		}); err != nil {
			return entity.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		id, err := deps.Mutator.Insert(ctx, port.InsertCommentDTO(input))
		if err != nil {
			return entity.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		comment, err := deps.Accessor.ByID(ctx, id)
		if err != nil {
			return entity.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		return comment, nil
	})
}
