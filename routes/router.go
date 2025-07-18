package routes

import (
	"short-url/internal/handlers"
	"short-url/internal/middleware"
	"short-url/internal/shortUrl"
	"short-url/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Setup(app *fiber.App) {
	app.Get("/", handlers.HelloWorld)
	api := app.Group("/api")
	api.Get("/health", handlers.HealthCheck)
	api.Get("/dashboard", monitor.New())
	auth := api.Group("/auth")
	auth.Post("/register", user.RegisterHandler)
	auth.Post("/login", user.LoginHandler)

	protected := api.Group("/url")
	protected.Use(middleware.JWTProtected())
	protected.Post("/create", shortUrl.CreateHandler)

	app.Get("/:code", shortUrl.RedirectHandler)
}
