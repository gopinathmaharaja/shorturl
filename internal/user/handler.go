package user

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterHandler(c *fiber.Ctx) error {
	var u User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}
	hash, _ := HashPassword(u.Password)
	u.Password = hash
	if err := CreateUser(&u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.JSON(fiber.Map{"message": "User created"})
}

func LoginHandler(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}
	user, err := FindOne(bson.M{"email": body["email"]})
	if err != nil || !CheckPassword(user.Password, body["password"]) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	token, _ := GenerateToken(user.ID)
	return c.JSON(fiber.Map{"token": token})
}
