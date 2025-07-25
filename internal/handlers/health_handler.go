package handlers

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "Server running successfully",
	})
}
func HelloWorld(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "Welcome to the Short URL service",
	})
}
