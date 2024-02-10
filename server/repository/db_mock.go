package repository

import (
	"context"
	"server/errs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errs.ErrDB
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errs.ErrDB
	}

	return client, nil
}