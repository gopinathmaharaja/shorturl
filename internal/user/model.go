package user

import (
	"time"
)

type User struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	Email          string    `bson:"email" json:"email"`
	Password       string    `bson:"password" json:"-"`
	APIKey         string    `bson:"api_key" json:"api_key"`
	TotalCount     int       `bson:"total_count" json:"total_count"`
	RemainingCount int       `bson:"remaining_count" json:"remaining_count"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}
