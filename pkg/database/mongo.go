package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection(ctx context.Context, username, password, uri string) (*mongo.Client, error) {
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})
	opts.ApplyURI(uri)

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := dbClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return dbClient, nil
}
