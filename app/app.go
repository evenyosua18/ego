package app

import (
	"github.com/evenyosua18/ego/config"
	"github.com/gofiber/fiber/v3"
)

type RouteHandler func(app *App)

type App struct {
	app *fiber.App
	cfg *Config
}

func Run(handler RouteHandler) {
	//  get config
	config.GetConfig()

	// map config
	cfg := BuildConfiguration()

	// route
	fiberApp := fiber.New()

	// db connection

	// create app
	app := &App{
		app: fiberApp,
		cfg: cfg,
	}

	// run handler
	handler(app)

	if err := fiberApp.Listen(cfg.ServiceConfig.Port); err != nil {
		panic(err)
	}
}
