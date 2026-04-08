package routes

import (
	"short-url/internal/handlers"
	"short-url/internal/middleware"
	"short-url/internal/shortUrl"
	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Setup(app *fiber.App, userCtrl *user.Controller, shortUrlCtrl *shortUrl.Controller, userRepo user.Repository) {
	app.Get("/", handlers.HelloWorld)
	api := app.Group("/api")
	api.Get("/health", handlers.HealthCheck)
	api.Get("/dashboard", monitor.New())

	auth := api.Group("/auth")
	auth.Use(middleware.RateLimit(5))
	auth.Post("/register", userCtrl.RegisterHandler)
	auth.Post("/login", userCtrl.LoginHandler)

	protected := api.Group("/url")
	protected.Use(middleware.JWTProtected(userRepo))
	protected.Use(middleware.RateLimit(1000))
	protected.Post("/create", shortUrlCtrl.CreateHandler)
	protected.Get("/analytics/:code", shortUrlCtrl.AnalyticsHandler) // New endpoint

	app.Get("/:code", middleware.RateLimit(1000), shortUrlCtrl.RedirectHandler)
}
