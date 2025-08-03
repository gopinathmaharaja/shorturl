package middleware

import (
	"os"
	"strings"

	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			apiKey := c.Get("X-API-KEY")
			user, err := user.FindOne(bson.M{"api_key": bson.M{"$eq": apiKey}})
			if err == nil && user != nil {
				c.Locals("userID", user.ID)
				return c.Next()
			}
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		_, err = user.FindOne(bson.M{"_id": bson.M{"$eq": userID}})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
