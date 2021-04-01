package helpers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GlobalMongoClient *mongo.Database

// ConnectDB : This is helper function to connect mongoDB
func ConnectDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}
	GlobalMongoClient = client.Database("movies_db")
	fmt.Println("Connected to MongoDB!")
}
