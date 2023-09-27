package mongodb

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	logger logging.Logger
	client *mongo.Client
}

func New(ctx context.Context, logger logging.Logger, connect string) (*DB, error) {
	const op = "db.New"
	var err error

	opts := options.Client().ApplyURI(connect)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &DB{
		logger: logger.With("system", "db", "dbType", "mongo"),
		client: client,
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	err := db.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("db.Close: %w", err)
	}

	return nil
}
