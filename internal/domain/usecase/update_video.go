package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/pkg/validation"
)

type UpdateVideoDeps struct {
	Accessor port.VideoAccessor
	Mutator  port.VideoMutator
}

func UpdateVideo(deps UpdateVideoDeps) port.UpdateVideo {
	return port.UsecaseFunc[port.UpdateVideoInput, entity.Video](func(
		ctx context.Context,
		input port.UpdateVideoInput,
	) (entity.Video, error) {
		const op = "usecase.UpdateVideo"

		if err := validation.Validate(func(v *validation.Validator) {
			if input.Data.Title != nil {
				vrule.Title(v, *input.Data.Title)
			}
			if input.Data.Description != nil {
				vrule.Description(v, *input.Data.Description)
			}
			if input.Data.ThumbnailPath != nil {
				vrule.ThumbnailPath(v, *input.Data.ThumbnailPath)
			}
			if input.Data.VideoPath != nil {
				vrule.VideoPath(v, *input.Data.VideoPath)
			}
		}); err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		if _, err := deps.Accessor.ByID(ctx, input.ID); err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		if err := deps.Mutator.Update(ctx, input.ID, port.UpdateVideoDTO(input.Data)); err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		newVideo, err := deps.Accessor.ByID(ctx, input.ID)
		if err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		return newVideo, nil
	})
}
