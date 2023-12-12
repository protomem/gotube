package repository

import (
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Comment interface{}

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
