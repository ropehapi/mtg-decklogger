package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const mongoOperationTimeout = 5 * time.Second

type MongoConnection struct {
	*mongo.Client
}

func ConnectMongo(mongoURI string) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoOperationTimeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		connection := &MongoConnection{Client: client}
		_ = connection.Close()
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	return &MongoConnection{Client: client}, nil
}

func (connection *MongoConnection) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoOperationTimeout)
	defer cancel()

	if err := connection.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect from MongoDB: %w", err)
	}

	return nil
}
