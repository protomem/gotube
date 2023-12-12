package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Comment = (*CommentImpl)(nil)

type CreateCommentDTO struct {
	Message  string
	AuthorID model.ID
	VideoID  model.ID
}

type (
	Comment interface {
		Get(ctx context.Context, id model.ID) (model.Comment, error)
		Create(ctx context.Context, dto CreateCommentDTO) (model.ID, error)
	}

	CommentImpl struct {
		logger logging.Logger
		mdb    *mongo.Client
	}
)

func NewComment(logger logging.Logger, mdb *mongo.Client) *CommentImpl {
	return &CommentImpl{
		logger: logger.With("repository", "comment", "repositoryType", "mongo"),
		mdb:    mdb,
	}
}

func (r *CommentImpl) Get(ctx context.Context, id model.ID) (model.Comment, error) {
	const op = "repository:Comment.Get"

	filter := bson.D{{Key: "commentId", Value: id.String()}}

	res := r.collection("comments").FindOne(ctx, filter)

	var doc commendDocuemnt
	if err := res.Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return model.Comment{}, fmt.Errorf("%s: %w", op, model.ErrCommentNotFound)
		}

		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return doc.Model(), nil
}

func (r *CommentImpl) Create(ctx context.Context, dto CreateCommentDTO) (model.ID, error) {
	const op = "repository:Comment.Create"

	doc, err := newCommentDocuemnt(dto)
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := r.collection("comments").InsertOne(ctx, doc); err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return doc.ModelID(), nil
}

func (r *CommentImpl) collection(name string) *mongo.Collection {
	return r.mdb.Database("gotubedb").Collection(name)
}

type commendDocuemnt struct {
	ID        primitive.ObjectID `bson:"_id"`
	CommentID string             `bson:"commentId"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	Message   string             `bson:"message"`
	AuthorID  string             `bson:"authorId"`
	VideoID   string             `bson:"videoId"`
}

func newCommentDocuemnt(dto CreateCommentDTO) (commendDocuemnt, error) {
	now := time.Now()

	commentID, err := uuid.NewRandom()
	if err != nil {
		return commendDocuemnt{}, err
	}

	return commendDocuemnt{
		ID:        primitive.NewObjectID(),
		CommentID: commentID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Message:   dto.Message,
		AuthorID:  dto.AuthorID.String(),
		VideoID:   dto.VideoID.String(),
	}, nil
}

// TODO: Handle panics
func (doc commendDocuemnt) ModelID() uuid.UUID {
	return uuid.MustParse(doc.CommentID)
}

func (doc commendDocuemnt) Model() model.Comment {
	return model.Comment{
		ID:        doc.ModelID(),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Message:   doc.Message,
		Author: model.User{
			ID: uuid.MustParse(doc.AuthorID),
		},
		VideoID: uuid.MustParse(doc.VideoID),
	}
}
