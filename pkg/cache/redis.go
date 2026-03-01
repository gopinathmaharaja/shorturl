package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		Client = nil
		return
	}

	log.Println("Successfully connected to Redis")
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return nil
	}
	return Client.Set(context.Background(), key, value, expiration).Err()
}

func Get(key string) (string, error) {
	if Client == nil {
		return "", nil
	}
	return Client.Get(context.Background(), key).Result()
}

func Delete(key string) error {
	if Client == nil {
		return nil
	}
	return Client.Del(context.Background(), key).Err()
}

func Close() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}
