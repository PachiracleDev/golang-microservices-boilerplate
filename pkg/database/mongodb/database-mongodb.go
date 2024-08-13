package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"moov/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase interface {
	GetCollection(collectionName string, databaseName *string, client *mongo.Client) *mongo.Collection
}

type Database struct {
	client       *mongo.Client
	databaseName string
}

func NewDatabase(config *config.Config) (*Database, error) {
	clientOptions := options.Client().ApplyURI(config.Database.Uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Database{
		client:       client,
		databaseName: config.Database.DatabaseName,
	}, nil
}

func (g *Database) GetCollection(collectionName string) *mongo.Collection {
	return g.client.Database(g.databaseName).Collection(collectionName)
}
