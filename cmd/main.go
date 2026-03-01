package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"short-url/internal/utils"
	"short-url/pkg/db"
	"short-url/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	app := fiber.New(fiber.Config{
		AppName: "Short URL Service v1.0",
	})

	// Initialize database
	db.Connect()
	defer db.Disconnect()

	// Initialize cron scheduler
	c := cron.New()
	utils.StartCleaningExpiredShortURLs(c)
	utils.StartMonthlyResetRemainingCount(c)
	c.Start()

	// Setup routes
	routes.Setup(app)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		c.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
