package mongo

import (
	"context"
	"log"
	"one1-be-chal/internal/adapters/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
	url string
}

func New(ctx context.Context, config *config.UserDB) (*DB, error) {
	clientOptions := options.Client().ApplyURI(config.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pinging MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")

	return &DB{
		Client: client,
		url:    config.URI,
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	err := db.Client.Disconnect(ctx)
	if err != nil {
		log.Println("Error disconnecting from MongoDB")
		return err
	}
	log.Println("Disconnected from MongoDB")
	return nil
}
