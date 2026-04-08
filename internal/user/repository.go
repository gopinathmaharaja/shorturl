package user

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(u *User) error
	UpdateMany(filter, update bson.M) (*mongo.UpdateResult, error)
	UpdateOne(filter, update bson.M) (*mongo.UpdateResult, error)
	FindOne(filter bson.M) (*User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &userRepository{collection: collection}
}

func (r *userRepository) CreateUser(u *User) error {
	log.Printf("[USER-REPO] Creating user with email: %s", u.Email)
	result, err := r.collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Printf("[USER-REPO] ERROR creating user %s: %v", u.Email, err)
		return err
	}
	u.ID = result.InsertedID.(primitive.ObjectID).Hex() // Assuming MongoDB uses ObjectId
	log.Printf("[USER-REPO] User created successfully with ID: %v, Email: %s", result.InsertedID, u.Email)
	return nil
}

func (r *userRepository) UpdateMany(filter, update bson.M) (*mongo.UpdateResult, error) {
	log.Printf("[USER-REPO] Updating multiple users. Filter: %+v", filter)
	result, err := r.collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Printf("[USER-REPO] ERROR updating multiple users: %v", err)
		return nil, err
	}
	log.Printf("[USER-REPO] Updated %d users", result.ModifiedCount)
	return result, nil
}

func (r *userRepository) UpdateOne(filter, update bson.M) (*mongo.UpdateResult, error) {
	log.Printf("[USER-REPO] Updating single user. Filter: %+v", filter)
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("[USER-REPO] ERROR updating user: %v", err)
		return nil, err
	}
	log.Printf("[USER-REPO] Update result - Matched: %d, Modified: %d", result.MatchedCount, result.ModifiedCount)
	return result, nil
}

func (r *userRepository) FindOne(filter bson.M) (*User, error) {
	log.Printf("[USER-REPO] Finding user with filter: %+v", filter)
	var user User
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
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
