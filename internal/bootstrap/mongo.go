package bootstrap

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoOptions struct {
	URI string
}

func Mongo(ctx context.Context, opts MongoOptions) (*mongo.Client, error) {
	const op = "bootstrap.Mongo"
	var err error

	mopts := options.Client().ApplyURI(opts.URI)
	client, err := mongo.Connect(ctx, mopts)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return client, nil
}
