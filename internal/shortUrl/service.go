package shortUrl

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateShortURL(original string, userID string) *ShortURL {
	return &ShortURL{
		Original:  original,
		ShortCode: generateCode(6),
		CreatedBy: userID,
	}
}
