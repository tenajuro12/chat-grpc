package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func InitMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017/")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("MongoDB client creation error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	fmt.Println("✅ Connected to MongoDB")
	return client.Database("chat_app")
}
