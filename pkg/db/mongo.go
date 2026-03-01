package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var client *mongo.Client

func Connect() {
	log.Println("Starting MongoDB connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}
	log.Println("MongoDB URI loaded.")

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("MongoDB client created.")

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}
	log.Println("Successfully connected and pinged MongoDB.")

	DB = client.Database("shorturl")
	log.Println("MongoDB database selected: shorturl")

	// Create indexes
	createIndexes()
}

func Disconnect() {
	if client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("MongoDB connection closed successfully")
	}
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Users collection indexes
	usersCollection := DB.Collection("users")
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := usersCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Printf("Error creating email index: %v", err)
	}

	// Short URLs collection indexes
	urlsCollection := DB.Collection("shorturls")

	// Index on short_code for fast lookups
	codeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "short_code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = urlsCollection.Indexes().CreateOne(ctx, codeIndex)
	if err != nil {
		log.Printf("Error creating short_code index: %v", err)
	}

	// Index on expire_at for cleanup
	expireIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "expire_at", Value: 1}},
	}
	_, err = urlsCollection.Indexes().CreateOne(ctx, expireIndex)
	if err != nil {
		log.Printf("Error creating expire_at index: %v", err)
	}

	// Compound index on created_by and created_at
	userDateIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_by", Value: 1},
			{Key: "created_at", Value: -1},
		},
	}
	_, err = urlsCollection.Indexes().CreateOne(ctx, userDateIndex)
	if err != nil {
		log.Printf("Error creating compound index: %v", err)
	}

	log.Println("Database indexes created successfully")
}
