package mongo

import (
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ repository.Comment = (*CommentRepository)(nil)

type CommentRepository struct {
	logger logging.Logger
	client *mongo.Client
}

func NewCommentRepository(logger logging.Logger, client *mongo.Client) *CommentRepository {
	return &CommentRepository{
		logger: logger.With("repository", "comment", "repositoryType", "mongo"),
		client: client,
	}
}
