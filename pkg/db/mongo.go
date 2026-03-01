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
	log.Println("[DATABASE] ========================================")
	log.Println("[DATABASE] Starting MongoDB connection...")
	log.Println("[DATABASE] ========================================")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("[DATABASE] FATAL - MONGO_URI environment variable not set")
	}
	log.Printf("[DATABASE] MongoDB URI loaded: %s", mongoURI)

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	log.Println("[DATABASE] Connection pool configured: Min=10, Max=100")

	var err error
	log.Println("[DATABASE] Attempting to connect to MongoDB...")
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("[DATABASE] FATAL - Failed to connect to MongoDB: %v", err)
	}
	log.Println("[DATABASE] MongoDB client created successfully")

	// Ping the database to verify connection
	log.Println("[DATABASE] Pinging MongoDB to verify connection...")
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("[DATABASE] FATAL - Could not ping MongoDB: %v", err)
	}
	log.Println("[DATABASE] MongoDB ping successful - connection verified")

	DB = client.Database("shorturl")
	log.Println("[DATABASE] MongoDB database selected: shorturl")

	// Create indexes
	log.Println("[DATABASE] Creating database indexes...")
	createIndexes()

	log.Println("[DATABASE] ========================================")
	log.Println("[DATABASE] MongoDB connection READY")
	log.Println("[DATABASE] ========================================")
}

func Disconnect() {
	if client == nil {
		log.Println("[DATABASE] WARNING - No MongoDB client to disconnect")
		return
	}

	log.Println("[DATABASE] Starting MongoDB disconnection...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("[DATABASE] ERROR disconnecting MongoDB: %v", err)
	} else {
		log.Println("[DATABASE] MongoDB connection closed successfully")
	}
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Users collection indexes
	usersCollection := DB.Collection("users")
	log.Println("[DATABASE-INDEX] Creating indexes for 'users' collection...")

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := usersCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Printf("[DATABASE-INDEX] WARNING - Error creating email index: %v", err)
	} else {
		log.Println("[DATABASE-INDEX] ✓ Email index created (unique)")
	}

	// Short URLs collection indexes
	urlsCollection := DB.Collection("shorturls")
	log.Println("[DATABASE-INDEX] Creating indexes for 'shorturls' collection...")

	// Index on short_code for fast lookups
	codeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "short_code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = urlsCollection.Indexes().CreateOne(ctx, codeIndex)
	if err != nil {
		log.Printf("[DATABASE-INDEX] WARNING - Error creating short_code index: %v", err)
	} else {
		log.Println("[DATABASE-INDEX] ✓ Short code index created (unique)")
	}

	// Index on expire_at for cleanup
	expireIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "expire_at", Value: 1}},
	}
	_, err = urlsCollection.Indexes().CreateOne(ctx, expireIndex)
	if err != nil {
		log.Printf("[DATABASE-INDEX] WARNING - Error creating expire_at index: %v", err)
	} else {
		log.Println("[DATABASE-INDEX] ✓ Expiration date index created")
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
		log.Printf("[DATABASE-INDEX] WARNING - Error creating compound index: %v", err)
	} else {
		log.Println("[DATABASE-INDEX] ✓ Compound index created (created_by, created_at)")
	}

	log.Println("[DATABASE-INDEX] ========================================")
	log.Println("[DATABASE-INDEX] All database indexes created successfully")
	log.Println("[DATABASE-INDEX] ========================================")
}
