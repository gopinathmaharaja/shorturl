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
	valid := err == nil
	if !valid {
		log.Printf("[USER-VALIDATION] Invalid email format: %s", email)
	}
	return valid
}

// ValidatePassword checks password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		log.Printf("[USER-VALIDATION] Password too short (length: %d)", len(password))
		return false, "Password must be at least 8 characters"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		log.Printf("[USER-VALIDATION] Password missing required complexity. Upper: %v, Lower: %v, Number: %v", hasUpper, hasLower, hasNumber)
		return false, "Password must contain uppercase, lowercase, and numbers"
	}

	return true, ""
}

func RegisterHandler(c *fiber.Ctx) error {
	clientIP := c.IP()
	log.Printf("[USER-REGISTER] New registration request from IP: %s", clientIP)

	var u User
	if err := c.BodyParser(&u); err != nil {
		log.Printf("[USER-REGISTER] ERROR parsing body from IP %s: %v", clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("[USER-REGISTER] Attempting to register email: %s from IP: %s", u.Email, clientIP)

	// Validate email
	if !ValidateEmail(u.Email) {
		log.Printf("[USER-REGISTER] FAILED - Invalid email format: %s from IP: %s", u.Email, clientIP)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	// Validate password strength
	if valid, msg := ValidatePassword(u.Password); !valid {
		log.Printf("[USER-REGISTER] FAILED - Password validation failed for %s from IP: %s - Reason: %s", u.Email, clientIP, msg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	// Check if user already exists
	_, err := FindOne(bson.M{"email": u.Email})
	if err == nil {
		log.Printf("[USER-REGISTER] FAILED - Email already exists: %s from IP: %s", u.Email, clientIP)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("[USER-REGISTER] ERROR checking user existence for %s from IP %s: %v", u.Email, clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	hash, err := HashPassword(u.Password)
	if err != nil {
		log.Printf("[USER-REGISTER] ERROR hashing password for %s from IP %s: %v", u.Email, clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing password"})
	}

	u.Password = hash
	u.TotalCount = 10
	u.RemainingCount = 10

	if err := CreateUser(&u); err != nil {
		log.Printf("[USER-REGISTER] ERROR creating user %s from IP %s: %v", u.Email, clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	log.Printf("[USER-REGISTER] SUCCESS - User registered: %s from IP: %s", u.Email, clientIP)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func LoginHandler(c *fiber.Ctx) error {
	clientIP := c.IP()
	log.Printf("[USER-LOGIN] New login request from IP: %s", clientIP)

	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		log.Printf("[USER-LOGIN] ERROR parsing body from IP %s: %v", clientIP, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	email := body["email"]
	log.Printf("[USER-LOGIN] Attempting login for email: %s from IP: %s", email, clientIP)

	// Validate input
	if email == "" || body["password"] == "" {
		log.Printf("[USER-LOGIN] FAILED - Missing credentials from IP: %s", clientIP)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password required"})
	}

	user, err := FindOne(bson.M{"email": email})
	if err == mongo.ErrNoDocuments {
		log.Printf("[USER-LOGIN] FAILED - User not found: %s from IP: %s", email, clientIP)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if err != nil {
		log.Printf("[USER-LOGIN] ERROR finding user %s from IP %s: %v", email, clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if !CheckPassword(user.Password, body["password"]) {
		log.Printf("[USER-LOGIN] FAILED - Invalid password for user: %s from IP: %s", email, clientIP)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		log.Printf("[USER-LOGIN] ERROR generating token for user %s from IP %s: %v", email, clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	log.Printf("[USER-LOGIN] SUCCESS - User logged in: %s from IP: %s", email, clientIP)
	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
