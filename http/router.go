package http

import (
	"context" // Added for ShutdownWithContext
	"fmt"

	"github.com/evenyosua18/ego/http/middleware"
	"github.com/evenyosua18/ego/tracer"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	html "github.com/gofiber/template/html/v2"
)

type Router struct {
	app fiber.Router
}

func NewRouter(cfg RouteConfig) *Router {
	// set fiber configs
	var fiberConfigs []fiber.Config

	// set html engine to fiber configs
	if cfg.HtmlPath != "" {
		// register html
		engine := html.New(cfg.HtmlPath, ".html")

		// set fiber configs
		fiberConfigs = append(fiberConfigs, fiber.Config{
			Views: engine,
		})
	}

	// route
	fiberApp := fiber.New(fiberConfigs...)

	// set favicon route if implement html engine
	if cfg.HtmlPath != "" {
		fiberApp.Get("/favicon.ico", func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNoContent)
		})
	}

	// set middleware
	fiberApp.Use(middleware.PanicHandler(), middleware.LogHandler())

	// set rate limiter middleware
	if cfg.MaxLimit != 0 {
		fiberApp.Use(middleware.RateLimiter(cfg.MaxLimit))
	}

	// set cors
	corsConfig := cors.ConfigDefault

	if len(cfg.CORS.AllowOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.CORS.AllowOrigins
	}

	if len(cfg.CORS.AllowMethods) > 0 {
		corsConfig.AllowMethods = cfg.CORS.AllowMethods
	}

	if len(cfg.CORS.AllowHeaders) > 0 {
		corsConfig.AllowHeaders = cfg.CORS.AllowHeaders
	}

	if cfg.CORS.AllowCredentials {
		corsConfig.AllowCredentials = cfg.CORS.AllowCredentials
	}

	fiberApp.Use(cors.New(corsConfig))

	// set router
	router := &Router{
		app: fiberApp,
	}

	// register routes
	RegisterRoutes(router.Group(cfg.MainPrefix))

	// print all registered routes
	if cfg.ShowRegisteredRoute {
		fmt.Println("registered routes")
		for _, routes := range fiberApp.Stack() {
			for _, route := range routes {
				if route.Path != "/" {
					fmt.Printf("%-6s %s\n", route.Method, route.Path)
				}
			}
		}
	}

	return router
}

func (r *Router) Listen(port string) error {
	return r.app.(*fiber.App).Listen(port)
}

func (r *Router) Shutdown() error {
	return r.app.(*fiber.App).Shutdown()
}

func (r *Router) ShutdownWithContext(ctx context.Context) error {
	return r.app.(*fiber.App).ShutdownWithContext(ctx)
}

func (r *Router) wrap(handler RouteHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		// setup context
		fiberCtx := fiberContext{ctx: c}

		// get http data
		data, opName := fiberCtx.HttpData()

		// create tracer
		sp := tracer.StartSpan(c.Context(), opName, tracer.WithAttributes(data))
		defer sp.End()

		return handler(&fiberCtx)
	}
}

func (r *Router) Get(path string, h RouteHandler, opts ...RouterFuncOption) {
	r.app.Get(path, r.wrap(h))
}

func (r *Router) Post(path string, h RouteHandler, opts ...RouterFuncOption) {
	r.app.Post(path, r.wrap(h))
}

func (r *Router) Put(path string, h RouteHandler, opts ...RouterFuncOption) {
	r.app.Put(path, r.wrap(h))
}

func (r *Router) Delete(path string, h RouteHandler, opts ...RouterFuncOption) {
	r.app.Delete(path, r.wrap(h))
}

func (r *Router) Patch(path string, h RouteHandler, opts ...RouterFuncOption) {
	r.app.Patch(path, r.wrap(h))
}

func (r *Router) Group(prefix string) IHttpRouter {
	return &Router{app: r.app.Group(prefix)}
}
