package http

import "github.com/gofiber/fiber/v3"

type Router struct {
	app fiber.Router
}

func NewRouter() *Router {
	// route
	fiberApp := fiber.New()

	// set router
	router := &Router{
		app: fiberApp,
	}

	// register routes
	RegisterRoutes(router)

	return router
}

func (r *Router) Listen(port string) error {
	return r.app.(*fiber.App).Listen(port)
}

func (r *Router) wrap(handler RouteHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		return handler(&fiberContext{ctx: c})
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
