package middleware

import (
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func LogHandler() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} - ${method} ${path} [${status}] ${latency}\n",
		TimeFormat: "2006/01/02 15:04:05",
	})
}
