package shortUrl

type ShortURL struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	Original  string `bson:"original" json:"original"`
	ShortCode string `bson:"short_code" json:"short_code"`
	CreatedBy string `bson:"created_by" json:"created_by"`
}
