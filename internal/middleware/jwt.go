package middleware

import (
	"log"
	"os"
	"strings"

	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()
		log.Printf("[JWT-MIDDLEWARE] Authentication attempt from IP: %s", clientIP)

		auth := c.Get("Authorization")
		if auth == "" {
			log.Printf("[JWT-MIDDLEWARE] No Authorization header from IP: %s, checking API key", clientIP)
			apiKey := c.Get("X-API-KEY")
			if apiKey == "" {
				log.Printf("[JWT-MIDDLEWARE] FAILED - No Authorization header or API key from IP: %s", clientIP)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
			}

			userData, err := user.FindOne(bson.M{"api_key": apiKey})
			if err == nil && userData != nil {
				log.Printf("[JWT-MIDDLEWARE] SUCCESS - API key authenticated for user: %s from IP: %s", userData.ID, clientIP)
				c.Locals("userID", userData.ID)
				return c.Next()
			}
			log.Printf("[JWT-MIDDLEWARE] FAILED - API key authentication failed from IP: %s: %v", clientIP, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		log.Printf("[JWT-MIDDLEWARE] Validating JWT token from IP: %s", clientIP)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Printf("[JWT-MIDDLEWARE] FAILED - JWT parse error from IP %s: %v", clientIP, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		if !token.Valid {
			log.Printf("[JWT-MIDDLEWARE] FAILED - Invalid token from IP: %s", clientIP)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("[JWT-MIDDLEWARE] FAILED - Could not parse claims from IP: %s", clientIP)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		userID, ok := claims["id"].(string)
		if !ok {
			log.Printf("[JWT-MIDDLEWARE] FAILED - Invalid user ID in token claims from IP: %s", clientIP)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		userData, err := user.FindOne(bson.M{"_id": userID})
		if err != nil {
			log.Printf("[JWT-MIDDLEWARE] FAILED - User not found (ID: %s) from IP: %s, error: %v", userID, clientIP, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		log.Printf("[JWT-MIDDLEWARE] SUCCESS - JWT authenticated for user: %s (email: %s) from IP: %s", userData.ID, userData.Email, clientIP)
		c.Locals("userID", userData.ID)
		return c.Next()
	}
}
