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

func (acc *VideoAccessor) AllWherePublic(ctx context.Context, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectWherePublic(ctx, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
}

func (acc *VideoAccessor) AllWherePublicAndSortByViews(ctx context.Context, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectWherePublicAndSortByViews(ctx, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
}

func (acc *VideoAccessor) AllByAuthorIDAndWherePublic(ctx context.Context, authorID entity.ID, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectByAuthorIDWherePublic(ctx, authorID, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
}

func (acc *VideoAccessor) AllByAuthorID(ctx context.Context, author entity.ID, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectByAuthorID(ctx, author, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
}

func (acc *VideoAccessor) AllByLikeTitleAndWherePublic(ctx context.Context, likeTitle string, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectByLikeTitleAndWherePublic(ctx, likeTitle, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
}

func (acc *VideoAccessor) AllByLikeTitleAndByAuthorIDAndWherePublic(ctx context.Context, likeTitle string, author entity.ID, opts port.FindOptions) ([]entity.Video, error) {
	videos, err := acc.videoDao.SelectByLikeTitleAndByAuthorIDAndWherePublic(ctx, likeTitle, author, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	authorIDs := make([]entity.ID, 0, len(videos))
	for _, video := range videos {
		authorIDs = append(authorIDs, video.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Video{}, nil
		}

		return []entity.Video{}, err
	}

	return mapVideoEntriesAndUserEntriesToVideoEntities(videos, authors), nil
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

func mapVideoEntriesAndUserEntriesToVideoEntities(videos []database.VideoEntry, users []database.UserEntry) []entity.Video {
	entities := make([]entity.Video, 0, len(videos))
	for _, video := range videos {
		for _, user := range users {
			if video.AuthorID == user.ID {
				entities = append(entities, mapVideoEntryAndUserEntryToVideoEntity(video, user))
				break
			}
		}
	}
	return entities
}
