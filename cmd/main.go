package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"short-url/pkg/db"
	"short-url/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	db.Connect()

	routes.Setup(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
