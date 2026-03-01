package user

import (
	"context"
	"log"

	"short-url/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return db.DB.Collection("users")
}

func CreateUser(u *User) error {
	log.Printf("[USER-REPO] Creating user with email: %s", u.Email)
	result, err := getCollection().InsertOne(context.TODO(), u)
	if err != nil {
		log.Printf("[USER-REPO] ERROR creating user %s: %v", u.Email, err)
		return err
	}
	log.Printf("[USER-REPO] User created successfully with ID: %v, Email: %s", result.InsertedID, u.Email)
	return nil
}

func UpdateMany(filter, update bson.M) (*mongo.UpdateResult, error) {
	log.Printf("[USER-REPO] Updating multiple users. Filter: %+v", filter)
	result, err := getCollection().UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Printf("[USER-REPO] ERROR updating multiple users: %v", err)
		return nil, err
	}
	log.Printf("[USER-REPO] Updated %d users", result.ModifiedCount)
	return result, nil
}

func UpdateOne(filter, update bson.M) (*mongo.UpdateResult, error) {
	log.Printf("[USER-REPO] Updating single user. Filter: %+v", filter)
	result, err := getCollection().UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("[USER-REPO] ERROR updating user: %v", err)
		return nil, err
	}
	log.Printf("[USER-REPO] Update result - Matched: %d, Modified: %d", result.MatchedCount, result.ModifiedCount)
	return result, nil
}

func FindOne(filter bson.M) (*User, error) {
	log.Printf("[USER-REPO] Finding user with filter: %+v", filter)
	var user User
	err := getCollection().FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Printf("[USER-REPO] User not found with filter: %+v", filter)
		return nil, err
	}
	if err != nil {
		log.Printf("[USER-REPO] ERROR finding user: %v", err)
		return nil, err
	}
	log.Printf("[USER-REPO] User found: ID=%s, Email=%s", user.ID, user.Email)
	return &user, nil
}
