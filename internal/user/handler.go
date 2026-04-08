package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) RegisterHandler(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()
	log.Printf("[USER-REGISTER] New registration request from IP: %s", clientIP)

	var u User
	if err := ctx.BodyParser(&u); err != nil {
		log.Printf("[USER-REGISTER] ERROR parsing body from IP %s: %v", clientIP, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("[USER-REGISTER] Attempting to register email: %s from IP: %s", u.Email, clientIP)

	if err := c.svc.Register(&u); err != nil {
		log.Printf("[USER-REGISTER] FAILED - %v", err)
		status := fiber.StatusBadRequest
		if err.Error() == "email already registered" {
			status = fiber.StatusConflict
		} else if err.Error() == "database error" || err.Error() == "error processing password" {
			status = fiber.StatusInternalServerError
		}
		return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[USER-REGISTER] SUCCESS - User registered: %s from IP: %s", u.Email, clientIP)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func (c *Controller) LoginHandler(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()
	log.Printf("[USER-LOGIN] New login request from IP: %s", clientIP)

	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil {
		log.Printf("[USER-LOGIN] ERROR parsing body from IP %s: %v", clientIP, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	email := body["email"]
	log.Printf("[USER-LOGIN] Attempting login for email: %s from IP: %s", email, clientIP)

	token, u, err := c.svc.Login(email, body["password"])
	if err != nil {
		log.Printf("[USER-LOGIN] FAILED - %v", err)
		status := fiber.StatusUnauthorized
		if err.Error() == "email and password required" {
			status = fiber.StatusBadRequest
		} else if err.Error() == "database error" || err.Error() == "could not generate token" {
			status = fiber.StatusInternalServerError
		}
		return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[USER-LOGIN] SUCCESS - User logged in: %s from IP: %s", email, clientIP)
	return ctx.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    u.ID,
			"email": u.Email,
		},
	})
}
