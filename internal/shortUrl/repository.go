package shortUrl

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateShortURL(doc *ShortURL) error
	DeleteShortURL(filter bson.M) error
	FindByCode(code string) (*ShortURL, error)
}

type shortUrlRepository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &shortUrlRepository{collection: collection}
}

func (r *shortUrlRepository) CreateShortURL(doc *ShortURL) error {
	log.Printf("[SHORTURL-REPO] Creating short URL: Code=%s, Original=%s, CreatedBy=%s", doc.ShortCode, doc.Original, doc.CreatedBy)
	result, err := r.collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Printf("[SHORTURL-REPO] ERROR creating short URL %s: %v", doc.ShortCode, err)
		return err
	}
	doc.ID = result.InsertedID.(primitive.ObjectID).Hex()
	log.Printf("[SHORTURL-REPO] Short URL created successfully. ID: %v, Code: %s", result.InsertedID, doc.ShortCode)
	return nil
}

func (r *shortUrlRepository) DeleteShortURL(filter bson.M) error {
	log.Printf("[SHORTURL-REPO] Deleting short URL with filter: %+v", filter)
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("[SHORTURL-REPO] ERROR deleting short URL: %v", err)
		return err
	}
	log.Printf("[SHORTURL-REPO] Short URL deleted. Deleted count: %d", result.DeletedCount)
	return nil
}

func (r *shortUrlRepository) FindByCode(code string) (*ShortURL, error) {
	log.Printf("[SHORTURL-REPO] Finding short URL with code: %s", code)
	var url ShortURL
	err := r.collection.FindOne(context.TODO(), bson.M{"short_code": code}).Decode(&url)
	if err == mongo.ErrNoDocuments {
		log.Printf("[SHORTURL-REPO] Short URL not found with code: %s", code)
		return nil, err
	}
	if err != nil {
		log.Printf("[SHORTURL-REPO] ERROR finding short URL by code %s: %v", code, err)
		return nil, err
	}
	log.Printf("[SHORTURL-REPO] Short URL found. Code: %s, Original: %s", url.ShortCode, url.Original)
	return &url, nil
}
