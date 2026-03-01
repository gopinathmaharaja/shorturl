package shortUrl

import (
	"log"
	"net/url"
	"time"

	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateHandler(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	// Validate original URL
	originalURL := body["original"]
	if originalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Original URL is required"})
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL format"})
	}

	userID := c.Locals("userID").(string)

	// Check user's remaining URL creation count
	count, err := GetUserShortURLCount(userID)
	if err != nil {
		log.Printf("Error retrieving user count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user count"})
	}
	if count <= 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL creation limit reached"})
	}

	short := GenerateShortURL(originalURL, userID)

	if err := CreateShortURL(short); err != nil {
		log.Printf("Error creating short URL: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save short URL"})
	}

	// Decrement user's count after successful creation
	if err := DecrementUserShortURLCount(userID); err != nil {
		log.Printf("Error decrementing user count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user count"})
	}

	return c.JSON(short)
}

// GetUserShortURLCount returns the number of remaining short URLs the user can create.
func GetUserShortURLCount(userID string) (int, error) {
	u, err := user.FindOne(bson.M{"_id": userID})
	if err != nil {
		return 0, err
	}
	return u.RemainingCount, nil
}

// DecrementUserShortURLCount decrements the user's remaining short URL creation count.
func DecrementUserShortURLCount(userID string) error {
	_, err := user.UpdateOne(
		bson.M{"_id": userID},
		bson.M{
			"$inc": bson.M{
				"remaining_count": -1,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func RedirectHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	url, err := FindByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	// Check if URL has expired
	if !url.ExpireAt.IsZero() && url.ExpireAt.Before(time.Now()) {
		return c.Status(fiber.StatusGone).SendString("URL has expired")
	}

	return c.Redirect(url.Original, fiber.StatusMovedPermanently)
}
