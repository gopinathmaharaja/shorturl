package shortUrl

import (
	"log"
	"net/url"
	"time"

	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateHandler(c *fiber.Ctx) error {
	clientIP := c.IP()
	userID := c.Locals("userID").(string)
	log.Printf("[SHORTURL-CREATE] New short URL creation request from user: %s (IP: %s)", userID, clientIP)

	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		log.Printf("[SHORTURL-CREATE] ERROR parsing body for user %s (IP: %s): %v", userID, clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	// Validate original URL
	originalURL := body["original"]
	if originalURL == "" {
		log.Printf("[SHORTURL-CREATE] FAILED - Empty original URL for user %s (IP: %s)", userID, clientIP)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Original URL is required"})
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		log.Printf("[SHORTURL-CREATE] FAILED - Invalid URL format for user %s (IP: %s): %s", userID, clientIP, originalURL)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL format"})
	}

	log.Printf("[SHORTURL-CREATE] Validating user quota for user: %s", userID)

	// Check user's remaining URL creation count
	count, err := GetUserShortURLCount(userID)
	if err != nil {
		log.Printf("[SHORTURL-CREATE] ERROR retrieving user count for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user count"})
	}

	log.Printf("[SHORTURL-CREATE] User %s has %d remaining URLs", userID, count)

	if count <= 0 {
		log.Printf("[SHORTURL-CREATE] FAILED - User quota exhausted for user %s (IP: %s)", userID, clientIP)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL creation limit reached"})
	}

	short := GenerateShortURL(originalURL, userID)
	log.Printf("[SHORTURL-CREATE] Generated short code: %s for user %s", short.ShortCode, userID)

	if err := CreateShortURL(short); err != nil {
		log.Printf("[SHORTURL-CREATE] ERROR saving short URL %s for user %s: %v", short.ShortCode, userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save short URL"})
	}

	// Decrement user's count after successful creation
	if err := DecrementUserShortURLCount(userID); err != nil {
		log.Printf("[SHORTURL-CREATE] WARNING - Could not decrement user count for user %s: %v", userID, err)
		// Don't fail the request for this
	}

	log.Printf("[SHORTURL-CREATE] SUCCESS - Short URL created. Code: %s, Original: %s, User: %s (IP: %s)", short.ShortCode, originalURL, userID, clientIP)
	return c.JSON(short)
}

// GetUserShortURLCount returns the number of remaining short URLs the user can create.
func GetUserShortURLCount(userID string) (int, error) {
	log.Printf("[SHORTURL-SERVICE] Fetching URL count for user: %s", userID)
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("[SHORTURL-SERVICE] FAILED - Invalid user ID format: %s: %v", userID, err)
		return 0, err
	}
	u, err := user.FindOne(bson.M{"_id": userObjID})
	if err != nil {
		log.Printf("[SHORTURL-SERVICE] ERROR fetching user data for %s: %v", userID, err)
		return 0, err
	}
	log.Printf("[SHORTURL-SERVICE] User %s has remaining count: %d", userID, u.RemainingCount)
	return u.RemainingCount, nil
}

// DecrementUserShortURLCount decrements the user's remaining short URL creation count.
func DecrementUserShortURLCount(userID string) error {
	log.Printf("[SHORTURL-SERVICE] Decrementing URL count for user: %s", userID)
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
	if err != nil {
		log.Printf("[SHORTURL-SERVICE] ERROR decrementing count for user %s: %v", userID, err)
		return err
	}
	log.Printf("[SHORTURL-SERVICE] Successfully decremented count for user: %s", userID)
	return nil
}

func RedirectHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")
	referer := c.Get("Referer")

	log.Printf("[SHORTURL-REDIRECT] Redirect request for code: %s from IP: %s, UserAgent: %s, Referer: %s", code, clientIP, userAgent, referer)

	url, err := FindByCode(code)
	if err != nil {
		log.Printf("[SHORTURL-REDIRECT] FAILED - Short URL not found: %s (IP: %s)", code, clientIP)
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	log.Printf("[SHORTURL-REDIRECT] Found URL code: %s, Original: %s", code, url.Original)

	// Check if URL has expired
	if !url.ExpireAt.IsZero() && url.ExpireAt.Before(time.Now()) {
		log.Printf("[SHORTURL-REDIRECT] FAILED - URL expired: %s (expired at: %v, now: %v)", code, url.ExpireAt, time.Now())
		return c.Status(fiber.StatusGone).SendString("URL has expired")
	}

	log.Printf("[SHORTURL-REDIRECT] SUCCESS - Redirecting code %s to: %s (IP: %s)", code, url.Original, clientIP)
	return c.Redirect(url.Original, fiber.StatusMovedPermanently)
}
