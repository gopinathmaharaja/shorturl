package middleware

import (
	"log"
	"os"
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

// isRateLimitSkipped checks if rate limiting should be skipped for dev/test environments
func isRateLimitSkipped(env string) bool {
	return env == "dev" || env == "test"
}

// cleanOldRequests filters out requests older than the time window
func cleanOldRequests(times []time.Time, windowStart time.Time) []time.Time {
	var valid []time.Time
	for _, t := range times {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	return valid
}

// checkRateLimit validates if the request count exceeds the limit
func checkRateLimit(ip string, validRequests []time.Time, limit int) error {
	if len(validRequests) >= limit {
		log.Printf("[RATELIMIT] Rate limit exceeded for IP: %s (limit: %d/min, current: %d requests)", ip, limit, len(validRequests))
		return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
	}
	return nil
}

// RateLimit middleware - configurable requests per minute per IP
// Skipped in dev and test environments for load testing
func RateLimit(requestsPerMinute int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		env := os.Getenv("ENV")
		log.Printf("[RATELIMIT] DEBUG - ENV variable value: '%s'", env)

		if isRateLimitSkipped(env) {
			log.Printf("[RATELIMIT] Rate limiting disabled in %s environment", env)
			return c.Next()
		}

		ip := c.IP()
		now := time.Now()
		windowStart := now.Add(-time.Minute)

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		times, exists := limiter.store[ip]
		if !exists {
			log.Printf("[RATELIMIT] First request from IP: %s (limit: %d/min)", ip, requestsPerMinute)
			limiter.store[ip] = []time.Time{now}
			return c.Next()
		}

		validRequests := cleanOldRequests(times, windowStart)
		limiter.store[ip] = validRequests

		if err := checkRateLimit(ip, validRequests, requestsPerMinute); err != nil {
			return err
		}

		log.Printf("[RATELIMIT] Request from IP: %s (requests in window: %d/%d)", ip, len(validRequests), requestsPerMinute)
		limiter.store[ip] = append(validRequests, now)
		return c.Next()
	}
}
