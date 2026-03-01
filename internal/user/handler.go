package user

import (
	"log"
	"net/mail"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ValidateEmail checks if email is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return false, "Password must contain uppercase, lowercase, and numbers"
	}

	return true, ""
}

func RegisterHandler(c *fiber.Ctx) error {
	var u User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate email
	if !ValidateEmail(u.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	// Validate password strength
	if valid, msg := ValidatePassword(u.Password); !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	// Check if user already exists
	_, err := FindOne(bson.M{"email": u.Email})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking user existence: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	hash, err := HashPassword(u.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing password"})
	}

	u.Password = hash
	u.TotalCount = 10
	u.RemainingCount = 10

	if err := CreateUser(&u); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func LoginHandler(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate input
	if body["email"] == "" || body["password"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password required"})
	}

	user, err := FindOne(bson.M{"email": body["email"]})
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if !CheckPassword(user.Password, body["password"]) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
