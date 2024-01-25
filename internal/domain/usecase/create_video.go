package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/vrule"
	"github.com/protomem/gotube/pkg/validation"
)

type CreateVideoDeps struct {
	Accessor port.VideoAccessor
	Mutator  port.VideoMutator
}

func CreateVideo(deps CreateVideoDeps) port.CreateVideo {
	return port.UsecaseFunc[port.CreateVideoInput, entity.Video](func(
		ctx context.Context,
		input port.CreateVideoInput,
	) (entity.Video, error) {
		const op = "usecase.CreateVideo"

		if err := validation.Validate(func(v *validation.Validator) {
			vrule.Title(v, input.Title)
			if input.Description != nil {
				vrule.Description(v, *input.Description)
			}
			if input.ThumbnailPath != nil {
				vrule.ThumbnailPath(v, *input.ThumbnailPath)
			}
			if input.VideoPath != nil {
				vrule.VideoPath(v, *input.VideoPath)
			}
		}); err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		dto := port.InsertVideoDTO{
			Title:         input.Title,
			Description:   "",
			ThumbnailPath: "",
			VideoPath:     "",
			AuthorID:      input.AuthorID,
			Public:        false,
		}
		if input.Description != nil {
			dto.Description = *input.Description
		}
		if input.ThumbnailPath != nil {
			dto.ThumbnailPath = *input.ThumbnailPath
		}
		if input.VideoPath != nil {
			dto.VideoPath = *input.VideoPath
		}
		if input.Public != nil {
			dto.Public = *input.Public
		}

		id, err := deps.Mutator.Insert(ctx, dto)
		if err != nil {
			return entity.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		video, err := deps.Accessor.ByID(ctx, id)
		if err != nil {
			return entity.Video{}, err
		}

		return video, nil
	})
}
