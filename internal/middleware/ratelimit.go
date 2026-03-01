package middleware

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type rateLimitStore struct {
	mu    sync.RWMutex
	store map[string][]time.Time
}

var limiter = &rateLimitStore{
	store: make(map[string][]time.Time),
}

// RateLimit middleware - configurable requests per minute per IP
func RateLimit(requestsPerMinute int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()
		windowStart := now.Add(-time.Minute)

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		// Clean old requests
		if times, exists := limiter.store[ip]; exists {
			var valid []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					valid = append(valid, t)
				}
			}
			limiter.store[ip] = valid

			if len(valid) >= requestsPerMinute {
				log.Printf("[RATELIMIT] Rate limit exceeded for IP: %s (limit: %d/min, current: %d requests)", ip, requestsPerMinute, len(valid))
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			}

			log.Printf("[RATELIMIT] Request from IP: %s (requests in window: %d/%d)", ip, len(valid), requestsPerMinute)
		} else {
			log.Printf("[RATELIMIT] First request from IP: %s (limit: %d/min)", ip, requestsPerMinute)
		}

		limiter.store[ip] = append(limiter.store[ip], now)
		return c.Next()
	}
}
