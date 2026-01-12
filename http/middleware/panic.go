package middleware

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/evenyosua18/ego/code"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func PanicHandler() fiber.Handler {
	// get default recover config
	config := recover.ConfigDefault

	// set enable stack trace
	config.EnableStackTrace = true

	// set stack trace handler
	config.StackTraceHandler = func(c fiber.Ctx, e any) {
		// capture panic message
		panicMessage := ""

		switch v := e.(type) {
		case string:
			panicMessage = v
		case runtime.Error:
			panicMessage = v.Error()
		case error:
			panicMessage = v.Error()
		default:
			panicMessage = fmt.Sprintf("unkown panic: %v", v)
		}

		// TODO Capture message to sentry

		log.Println(panicMessage)
		fmt.Println(string(debug.Stack()))

		errCode := code.Get(code.PanicError)
		_ = c.Status(errCode.HttpCode).JSON(errCode)
	}

	return recover.New(config)
}
