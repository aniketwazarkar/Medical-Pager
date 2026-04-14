package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"medical-pager/utils"
)

var Client *mongo.Client

// Connect establishes a connection to MongoDB
func Connect() {
	uri := utils.GetEnv("MONGODB_URI", "mongodb://localhost:27017/medical-pager")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Output connected status
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	Client = client
}

// GetCollection returns a MongoDB collection reference
func GetCollection(collectionName string) *mongo.Collection {
	// Parse the db name from URI or use a default.
	// For production readiness, you might want to specify the DB name explicitly in an env var too,
	// but standard Mongo URI handling allows specifying the db in the URI.
	// We'll use "medical-pager" as default.
	return Client.Database("medical-pager").Collection(collectionName)
}
