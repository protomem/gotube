package adapter

import (
	"context"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/infra/database"
)

var _ port.VideoAccessor = (*VideoAccessor)(nil)

type VideoAccessor struct {
	videoDao *database.VideoDAO
	userDao  *database.UserDAO
}

func NewVideoAccessor(db *database.DB) *VideoAccessor {
	return &VideoAccessor{
		videoDao: db.VideoDAO(),
		userDao:  db.UserDAO(),
	}
}

func (acc *VideoAccessor) ByID(ctx context.Context, id entity.ID) (entity.Video, error) {
	video, err := acc.videoDao.GetByID(ctx, id)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.Video{}, entity.ErrVideoNotFound
		}
	}

	author, err := acc.userDao.GetByID(ctx, video.AuthorID)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.Video{}, entity.ErrVideoNotFound
		}

		return entity.Video{}, err
	}

	return mapVideoEntryAndUserEntryToVideoEntity(video, author), nil
}

var _ port.VideoMutator = (*VideoMutator)(nil)

type VideoMutator struct {
	dao *database.VideoDAO
}

func NewVideoMutator(db *database.DB) *VideoMutator {
	return &VideoMutator{db.VideoDAO()}
}

func (mut *VideoMutator) Insert(ctx context.Context, dto port.InsertVideoDTO) (entity.ID, error) {
	id, err := mut.dao.Insert(ctx, database.InsertVideoDTO(dto))
	if err != nil {
		if database.IsKeyConflict(err) {
			return entity.ID{}, entity.ErrVideoAlreadyExists
		}

		return entity.ID{}, err
	}

	return id, nil
}

func (mut *VideoMutator) Update(ctx context.Context, id entity.ID, dto port.UpdateVideoDTO) error {
	return mut.dao.Update(ctx, id, database.UpdateVideoDTO(dto))
}

func (mut *VideoMutator) Delete(ctx context.Context, id entity.ID) error {
	return mut.dao.Delete(ctx, id)
}

func mapVideoEntryAndUserEntryToVideoEntity(video database.VideoEntry, user database.UserEntry) entity.Video {
	return entity.Video{
		ID:            video.ID,
		CreatedAt:     video.CreatedAt,
		UpdatedAt:     video.UpdatedAt,
		Title:         video.Title,
		Description:   video.Description,
		ThumbnailPath: video.ThumbnailPath,
		VideoPath:     video.VideoPath,
		Author:        entity.User(user),
		Public:        video.Public,
		Views:         video.Views,
	}
}
