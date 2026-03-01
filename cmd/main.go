package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"short-url/internal/utils"
	"short-url/pkg/db"
	"short-url/pkg/logger"
	"short-url/routes"
)

func main() {
	// Initialize logger first
	logger.Init()
	logger.InfoLog.Println("========================================")
	logger.InfoLog.Println("Short URL Service Starting")
	logger.InfoLog.Println("========================================")

	err := godotenv.Load()
	if err != nil {
		logger.WarnLog.Println("Warning: .env file not found, using environment variables")
	} else {
		logger.InfoLog.Println(".env file loaded successfully")
	}

	logger.InfoLog.Printf("Environment: %s", os.Getenv("ENV"))

	app := fiber.New(fiber.Config{
		AppName: "Short URL Service v1.0",
	})

	// Initialize database
	logger.InfoLog.Println("Initializing database connection...")
	db.Connect()
	logger.InfoLog.Println("Database connection established and ready")
	defer func() {
		logger.InfoLog.Println("Closing database connection...")
		db.Disconnect()
	}()

	// Initialize cron scheduler
	logger.InfoLog.Println("Initializing cron scheduler...")
	c := cron.New()
	utils.StartCleaningExpiredShortURLs(c)
	logger.InfoLog.Println("Cron job registered: Cleaning expired short URLs")
	utils.StartMonthlyResetRemainingCount(c)
	logger.InfoLog.Println("Cron job registered: Monthly reset remaining count")
	c.Start()
	logger.InfoLog.Println("Cron scheduler started")

	// Setup routes
	logger.InfoLog.Println("Setting up routes...")
	routes.Setup(app)
	logger.InfoLog.Println("Routes configured successfully")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan

		logger.InfoLog.Printf("Received signal: %v", sig)
		logger.InfoLog.Println("========================================")
		logger.InfoLog.Println("Shutting down server gracefully...")
		logger.InfoLog.Println("========================================")

		c.Stop()
		logger.InfoLog.Println("Cron scheduler stopped")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.ErrorLog.Printf("Server forced to shutdown with error: %v", err)
		} else {
			logger.InfoLog.Println("Server shutdown completed successfully")
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		logger.WarnLog.Println("PORT not set, using default: 3000")
	}

	logger.InfoLog.Printf("Server listening on port %s", port)
	logger.InfoLog.Println("========================================")
	logger.InfoLog.Println("Short URL Service is now READY")
	logger.InfoLog.Println("========================================")

	if err := app.Listen(":" + port); err != nil {
		logger.ErrorLog.Printf("Error starting server: %v", err)
	}
}
