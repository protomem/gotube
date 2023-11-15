package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	_database           = "gotubedb"
	_commentCollections = "comments"
)

var _ repository.Comment = (*CommentRepository)(nil)

type CommentRepository struct {
	logger logging.Logger
	client *mongo.Client

	userRepo repository.User
}

func NewCommentRepository(logger logging.Logger, client *mongo.Client, userRepo repository.User) *CommentRepository {
	return &CommentRepository{
		logger:   logger.With("repository", "comment", "repositoryType", "mongo"),
		client:   client,
		userRepo: userRepo,
	}
}

func (repo *CommentRepository) FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Comment, error) {
	const op = "mongo.CommentRepository.FindByVideoID"
	var err error

	filter := bson.D{{Key: "videoId", Value: videoID.String()}}

	cursor, err := repo.client.
		Database(_database).
		Collection(_commentCollections).
		Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Comment{}, nil
		}

		return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var comments []model.Comment
	for cursor.Next(ctx) {
		var doc commentDocument
		err := cursor.Decode(&doc)
		if err != nil {
			return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		comment, err := doc.model()
		if err != nil {
			return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		comments = append(comments, comment)
	}

	if cursor.Err() != nil {
		return []model.Comment{}, fmt.Errorf("%s: %w", op, cursor.Err())
	}

	if len(comments) == 0 {
		return []model.Comment{}, nil
	}

	// TODO: ? ...
	{
		usersIDs := make([]uuid.UUID, 0, len(comments))
		for _, comment := range comments {
			usersIDs = append(usersIDs, comment.Author.ID)
		}

		users, err := repo.userRepo.Find(ctx, usersIDs)
		if err != nil {
			return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
		}

		for _, user := range users {
			for i, comment := range comments {
				if comment.Author.ID == user.ID {
					comments[i].Author = user
				}
			}
		}
	}

	return comments, nil
}

func (repo *CommentRepository) Get(ctx context.Context, id uuid.UUID) (model.Comment, error) {
	const op = "mongo.CommentRepository.Get"

	filter := bson.D{{Key: "commentId", Value: id.String()}}

	res := repo.client.
		Database(_database).
		Collection(_commentCollections).
		FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return model.Comment{}, fmt.Errorf("%s: %w", op, model.ErrCommentNotFound)
		}

		return model.Comment{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	var doc commentDocument
	err := res.Decode(&doc)
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment, err := doc.model()
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: ? ...
	{
		comment.Author, err = repo.userRepo.Get(ctx, comment.Author.ID)
		if err != nil {
			return model.Comment{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return comment, nil
}

func (repo *CommentRepository) Create(ctx context.Context, dto repository.CreateCommentDTO) (uuid.UUID, error) {
	const op = "mongo.CommentRepository.Create"
	var err error

	doc, err := newCommentDocument(dto)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = repo.client.
		Database(_database).
		Collection(_commentCollections).
		InsertOne(ctx, doc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	commentID, _ := uuid.Parse(doc.CommentID)

	return commentID, nil
}

type commentDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	CommentID string             `bson:"commentId"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`

	Content string `bson:"content"`

	AuthorID string `bson:"authorId"`
	VideoID  string `bson:"videoId"`
}

func newCommentDocument(dto repository.CreateCommentDTO) (commentDocument, error) {
	now := time.Now()

	commentID, err := uuid.NewRandom()
	if err != nil {
		return commentDocument{}, fmt.Errorf("generate comment id: %w", err)
	}

	return commentDocument{
		ID:        primitive.NewObjectID(),
		CommentID: commentID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Content:   dto.Content,
		AuthorID:  dto.AuthorID.String(),
		VideoID:   dto.VideoID.String(),
	}, nil
}

func (doc *commentDocument) model() (model.Comment, error) {
	var err error

	commentID, err := uuid.Parse(doc.CommentID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("invalid comment id: %w", err)
	}

	authorID, err := uuid.Parse(doc.AuthorID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("invalid author id: %w", err)
	}

	videoID, err := uuid.Parse(doc.VideoID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("invalid video id: %w", err)
	}

	return model.Comment{
		ID:        commentID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Content:   doc.Content,
		Author: model.User{
			ID: authorID,
		},
		VideoID: videoID,
	}, nil
}
