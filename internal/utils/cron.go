package utils

import (
	"log"
	"time"

	"short-url/internal/shortUrl"
	"short-url/internal/user"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func StartCleaningExpiredShortURLs(c *cron.Cron, repo shortUrl.Repository) {
	_, err := c.AddFunc("0 0 1 * *", func() {
		err := repo.DeleteShortURL(bson.M{
			"expire_at": bson.M{
				"$lt": time.Now(),
			},
		})
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
}
func StartMonthlyResetRemainingCount(c *cron.Cron, repo user.Repository) {
	_, err := c.AddFunc("0 0 1 * *", func() {
		_, err := repo.UpdateMany(bson.M{}, bson.M{
			"$set": bson.M{
				"remaining_count": "$total_count",
			},
		})
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
}
