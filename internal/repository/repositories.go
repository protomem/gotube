package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct{}

func New(logger logging.Logger, pdb *pgxpool.Pool, mdb *mongo.Client) *Repositories {
	return &Repositories{}
}
