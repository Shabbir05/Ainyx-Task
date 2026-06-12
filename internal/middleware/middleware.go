package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID injects a UUID into the X-Request-ID response header and
// stores it in the Fiber Locals map under the key "requestID".
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Re-use an existing request ID if the client already sent one.
		id := c.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set("X-Request-ID", id)
		c.Locals("requestID", id)
		return c.Next()
	}
}

// RequestLogger logs method, path, status code and duration using Zap
// after every request completes.
func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		requestID, _ := c.Locals("requestID").(string)

		log.Info("request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)

		return err
	}
}
