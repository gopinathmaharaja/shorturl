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
	short := GenerateShortURL(body["original"], userID)

	if err := CreateShortURL(short); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save short URL"})
	}

	return c.JSON(short)
}

func RedirectHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	url, err := FindByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}
	return c.Redirect(url.Original, fiber.StatusMovedPermanently)
}
