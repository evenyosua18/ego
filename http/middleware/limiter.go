package middleware

import (
	"github.com/evenyosua18/ego/code"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func RateLimiter(maxLimit int) fiber.Handler {
	// get default config
	config := limiter.ConfigDefault

	// set max
	config.Max = maxLimit

	// set limit response
	config.LimitReached = func(c fiber.Ctx) error {
		// get codes
		errCode := code.Get(code.RateLimitError)

		return c.Status(errCode.HttpCode).JSON(errCode)
	}

	// set key generator
	config.KeyGenerator = func(c fiber.Ctx) string {
		return c.IP()
	}

	return limiter.New(config)
}
