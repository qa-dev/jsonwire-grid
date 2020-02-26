package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type Factory struct {
}

func (f *Factory) Create(cfg config.Config) (pool.StorageInterface, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.Connection))
	if err != nil {
		panic("Database connection error: " + err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		err = errors.New("Database connection not established: " + err.Error())
		return nil, err
	}

	db := client.Database(cfg.DB.DbName)
	s := NewMongoStorage(db)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"key":     1,
			"address": 1,
		},
		Options: options.Index().SetUnique(true).SetName("key_address_unique"),
	}

	_, err = s.collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		err = errors.New("Create index error: " + err.Error())
		return nil, err
	}
	return s, nil
}
