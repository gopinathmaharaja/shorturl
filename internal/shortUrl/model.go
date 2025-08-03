package shortUrl

import (
	"time"
)

type ShortURL struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Original  string    `bson:"original" json:"original"`
	ShortCode string    `bson:"short_code" json:"short_code"`
	CreatedBy string    `bson:"created_by" json:"created_by"`
	ExpireAt  time.Time `bson:"expire_at" json:"expire_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
