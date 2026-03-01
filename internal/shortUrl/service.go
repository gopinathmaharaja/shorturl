package shortUrl

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func GenerateShortURL(original string, userID string) *ShortURL {
	return &ShortURL{
		Original:  original,
		ShortCode: generateCode(6),
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
