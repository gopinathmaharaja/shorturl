package shortUrl

import (
	"context"
	"short-url/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return db.DB.Collection("shorturls")
}

func CreateShortURL(doc *ShortURL) error {
	_, err := getCollection().InsertOne(context.TODO(), doc)
	return err
}

func FindByCode(code string) (*ShortURL, error) {
	var url ShortURL
	err := getCollection().FindOne(context.TODO(), bson.M{"short_code": code}).Decode(&url)
	return &url, err
}
