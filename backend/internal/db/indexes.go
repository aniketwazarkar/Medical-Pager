package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupIndexes creates the required database indexes for scalability and performance
func SetupIndexes() {
	ctx := context.Background()

	// Generic function to create an index
	createIndex := func(collectionName string, keys bson.D) {
		coll := GetCollection(collectionName)
		idxModel := mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetBackground(true),
		}
		if _, err := coll.Indexes().CreateOne(ctx, idxModel); err != nil {
			log.Printf("Warning: failed to create index on %s: %v", collectionName, err)
		}
	}

	// 1. users: index tenantId
	createIndex("users", bson.D{{Key: "tenantId", Value: 1}})
	createIndex("users", bson.D{{Key: "email", Value: 1}}) // Also index email for login speed

	// 2. messages: index tenantId, channelId, senderId, createdAt
	createIndex("messages", bson.D{{Key: "tenantId", Value: 1}})
	createIndex("messages", bson.D{{Key: "channelId", Value: 1}})
	createIndex("messages", bson.D{{Key: "senderId", Value: 1}})
	createIndex("messages", bson.D{{Key: "createdAt", Value: -1}})

	// 3. channels: index tenantId
	createIndex("channels", bson.D{{Key: "tenantId", Value: 1}})

	// 4. patients: index tenantId
	createIndex("patients", bson.D{{Key: "tenantId", Value: 1}})
	createIndex("patients", bson.D{{Key: "patientId", Value: 1}})

	// 5. audit_logs: index tenantId, createdAt
	createIndex("audit_logs", bson.D{{Key: "tenantId", Value: 1}})
	createIndex("audit_logs", bson.D{{Key: "createdAt", Value: -1}})

	log.Println("Database indexes initialized successfully.")
}
