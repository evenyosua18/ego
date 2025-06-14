package http

import (
	"fmt"
	"github.com/evenyosua18/ego/http/middleware"
	"github.com/evenyosua18/ego/tracer"
	"github.com/gofiber/fiber/v3"
)

type Router struct {
	app fiber.Router
}

func NewRouter(cfg RouteConfig) *Router {
	// route
	fiberApp := fiber.New()

	// set middleware
	fiberApp.Use(middleware.PanicHandler(), middleware.LogHandler())

	// set rate limiter middleware
	if cfg.MaxLimit != 0 {
		fiberApp.Use(middleware.RateLimiter(cfg.MaxLimit))
	}

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
