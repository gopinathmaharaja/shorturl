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

func UpdateMany(filter, update bson.M) (*mongo.UpdateResult, error) {
	return getCollection().UpdateMany(context.TODO(), filter, update)
}

func UpdateOne(filter, update bson.M) (*mongo.UpdateResult, error) {
	return getCollection().UpdateOne(context.TODO(), filter, update)
}

func FindOne(filter bson.M) (*User, error) {
	var user User
	err := getCollection().FindOne(context.TODO(), filter).Decode(&user)
	return &user, err
}
