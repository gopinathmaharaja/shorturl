package user

import (
	"context"
	"short-url/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return db.DB.Collection("users")
}

func CreateUser(u *User) error {
	_, err := getCollection().InsertOne(context.TODO(), u)
	return err
}

func FindByEmail(email string) (*User, error) {
	var user User
	err := getCollection().FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return &user, err
}
