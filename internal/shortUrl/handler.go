package shortUrl

import (
	"github.com/gofiber/fiber/v2"
)

func CreateHandler(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}
	userID := c.Locals("userID").(string)

	// Check user's remaining URL creation count
	count, err := GetUserShortURLCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user count"})
	}
	if count <= 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL creation limit reached"})
	}

	short := GenerateShortURL(body["original"], userID)

	if err := CreateShortURL(short); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save short URL"})
	}

	// Decrement user's count after successful creation
	if err := DecrementUserShortURLCount(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user count"})
	}

	return c.JSON(short)
}

// GetUserShortURLCount returns the number of remaining short URLs the user can create.
func GetUserShortURLCount(userID string) (int, error) {
	// TODO: Implement actual logic to fetch user's remaining count from database or cache.
	// For now, return a dummy value for demonstration.
	return 5, nil
}

// DecrementUserShortURLCount decrements the user's remaining short URL creation count.
func DecrementUserShortURLCount(userID string) error {
	// TODO: Implement actual logic to decrement user's count in database .
	// For now, return nil for demonstration.
	return nil
}

func RedirectHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	url, err := FindByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}
	return c.Redirect(url.Original, fiber.StatusMovedPermanently)
}
