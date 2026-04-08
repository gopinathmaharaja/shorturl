package shortUrl

import (
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	svc          Service
	analyticsRepo AnalyticsRepository
}

func NewController(svc Service, analyticsRepo AnalyticsRepository) *Controller {
	return &Controller{svc: svc, analyticsRepo: analyticsRepo}
}

func (c *Controller) CreateHandler(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()
	userID := ctx.Locals("userID").(string)
	log.Printf("[SHORTURL-CREATE] New short URL creation request from user: %s (IP: %s)", userID, clientIP)

	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil {
		log.Printf("[SHORTURL-CREATE] ERROR parsing body for user %s (IP: %s): %v", userID, clientIP, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	// Validate original URL
	originalURL := body["original"]
	if originalURL == "" {
		log.Printf("[SHORTURL-CREATE] FAILED - Empty original URL for user %s (IP: %s)", userID, clientIP)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Original URL is required"})
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		log.Printf("[SHORTURL-CREATE] FAILED - Invalid URL format for user %s (IP: %s): %s", userID, clientIP, originalURL)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL format"})
	}

	customAlias := body["custom_alias"] // May be empty, that's fine

	short, err := c.svc.CreateShortURL(originalURL, customAlias, userID)
	if err != nil {
		log.Printf("[SHORTURL-CREATE] FAILED for user %s: %v", userID, err)
		status := fiber.StatusInternalServerError
		if err.Error() == "url creation limit reached" {
			status = fiber.StatusForbidden
		} else if err.Error() == "custom alias already in use" {
			status = fiber.StatusConflict
		} else if err.Error() == "custom alias must be between 3 and 20 characters" {
			status = fiber.StatusBadRequest
		}
		return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[SHORTURL-CREATE] SUCCESS - Short URL created. Code: %s, Original: %s, User: %s (IP: %s)", short.ShortCode, originalURL, userID, clientIP)
	return ctx.JSON(short)
}

func (c *Controller) RedirectHandler(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	clientIP := ctx.IP()
	userAgent := ctx.Get("User-Agent")
	referer := ctx.Get("Referer")

	log.Printf("[SHORTURL-REDIRECT] Redirect request for code: %s from IP: %s, UserAgent: %s, Referer: %s", code, clientIP, userAgent, referer)

	url, err := c.svc.GetShortURL(code)
	if err != nil {
		log.Printf("[SHORTURL-REDIRECT] FAILED - Short URL not found: %s (IP: %s)", code, clientIP)
		return ctx.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	// Check if URL has expired
	if !url.ExpireAt.IsZero() && url.ExpireAt.Before(time.Now()) {
		log.Printf("[SHORTURL-REDIRECT] FAILED - URL expired: %s (expired at: %v, now: %v)", code, url.ExpireAt, time.Now())
		return ctx.Status(fiber.StatusGone).SendString("URL has expired")
	}

	// Track analytics asynchronously
	go func(sc string) {
		if err := c.analyticsRepo.TrackClick(sc); err != nil {
			log.Printf("[SHORTURL-REDIRECT-ANALYTICS] Failed to track click for %s: %v", sc, err)
		}
	}(code)

	log.Printf("[SHORTURL-REDIRECT] SUCCESS - Redirecting code %s to: %s (IP: %s)", code, url.Original, clientIP)
	return ctx.Redirect(url.Original, fiber.StatusMovedPermanently)
}

func (c *Controller) AnalyticsHandler(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	log.Printf("[SHORTURL-ANALYTICS] Request for analytics of code: %s", code)

	// In a complete implementation we might check if the user requesting owns the ShortURL,
	// but for now anyone authenticated can see analytics of a URL code.
	
	analyticsData, err := c.analyticsRepo.GetAnalytics(code)
	if err != nil {
		// It might be just 0 clicks. So return default struct.
		return ctx.JSON(&Analytics{
			ShortCode: code,
			Clicks: 0,
		})
	}

	return ctx.JSON(analyticsData)
}
