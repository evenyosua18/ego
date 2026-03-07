package http

import (
	"context" // Added for ShutdownWithContext
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/evenyosua18/ego/auth"
	"github.com/evenyosua18/ego/code"
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/http/middleware"
	"github.com/evenyosua18/ego/tracer"
	"github.com/gofiber/contrib/v3/swaggerui"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	html "github.com/gofiber/template/html/v2"
)

type Router struct {
	app                fiber.Router
	DisableAuthChecker bool
	routeGroup         string
}

func NewRouter(cfg RouteConfig) *Router {
	// set fiber configs
	fiberConfig := fiber.Config{
		Concurrency:  cfg.MaxConnection,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// set html engine to fiber configs
	if cfg.HtmlPath != "" {
		// register html
		engine := html.New(cfg.HtmlPath, ".html")
		fiberConfig.Views = engine
	}

	// route
	fiberApp := fiber.New(fiberConfig)

	// setup documentation
	if cfg.Doc.Path != "" {
		// fiber-swagger middleware
		fiberApp.Get("/docs", swaggerui.New(swaggerui.Config{
			BasePath: "/",
			FilePath: cfg.Doc.Path,
			Path:     "docs",
		}))

		// set swagger
		fiberApp.Get("/"+cfg.Doc.Path, func(c fiber.Ctx) error {
			// Returns the file located at ./docs/swagger.json
			return c.SendFile(cfg.Doc.Path)
		})
	}

	// set favicon route if implement html engine
	if cfg.HtmlPath != "" {
		fiberApp.Get("/favicon.ico", func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNoContent)
		})
	}

	// set middleware
	fiberApp.Use(middleware.PanicHandler(), middleware.LogHandler())

	// set rate limiter middleware
	if cfg.RateLimit.MaxLimit != 0 {
		fiberApp.Use(middleware.RateLimiter(cfg.RateLimit.MaxLimit))
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
		app:                fiberApp,
		DisableAuthChecker: cfg.DisableAuthChecker,
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

func (r *Router) ActiveConnections() int {
	if fiberApp, ok := r.app.(*fiber.App); ok {
		return int(fiberApp.Server().GetOpenConnectionsCount())
	}
	return 0
}

func (r *Router) Test(req *http.Request) (*http.Response, error) {
	if fiberApp, ok := r.app.(*fiber.App); ok {
		return fiberApp.Test(req)
	}
	return nil, fmt.Errorf("router is not a fiber app")
}

func (r *Router) extractWrap(h RouteHandler, opts []RouterFuncOption) (fiber.Handler, []fiber.Handler) {
	// apply options
	routeOpt := &RouterOptions{}
	for _, opt := range opts {
		opt(routeOpt)
	}

	// append a middleware to inject route roles into fiber local context
	var middlewares []fiber.Handler
	if r.DisableAuthChecker {
		middlewares = make([]fiber.Handler, 0, len(routeOpt.Middlewares)+1)
		if len(routeOpt.Roles) > 0 {
			middlewares = append(middlewares, func(c fiber.Ctx) error {
				c.Locals(localRouteRoles{}, routeOpt.Roles)
				return c.Next()
			})
		}

		middlewares = append(middlewares, routeOpt.Middlewares...)
	} else {
		middlewares = routeOpt.Middlewares
	}

	handler := func(c fiber.Ctx) error {
		// setup context
		fiberCtx := fiberContext{ctx: c}

		// auth validation
		if !r.DisableAuthChecker {
			// get authorization header
			authHeader := c.Get("Authorization")

			// check length
			if len(authHeader) == 0 {
				return code.Get(code.BadRequestError).SetMessage("header not found")
			}

			// split token
			splitToken := strings.Split(authHeader, "Bearer ")

			if len(splitToken) != 2 {
				return code.Get(code.BadRequestError).SetMessage("invalid authorization header")
			}

			// validate access token
			claims, err := auth.ValidateToken(c.Context(), splitToken[1])
			if err != nil {
				return code.Get(code.UnauthorizedError)
			}

			// validate access type
			if claims.AccessType != r.routeGroup {
				return code.Get(code.UnauthorizedError).SetMessage("invalid access type")
			}

			switch r.routeGroup {
			case "svc":
				// validate service eligibility
				if err = validateServiceToken(claims.AllowedServices); err != nil {
					return err
				}
				break
			case "public":
				// validate public eligibility
				if err = validatePublicToken(routeOpt.Roles, claims.Roles); err != nil {
					return err
				}
				break
			}

		}

		// get http data
		data, opName := fiberCtx.HttpData()

		// create tracer
		sp := tracer.StartSpan(c.Context(), opName, tracer.WithAttributes(data))
		defer sp.End()

		return h(&fiberCtx)
	}

	return handler, middlewares
}

func (r *Router) Get(path string, h RouteHandler, opts ...RouterFuncOption) {
	handler, middlewares := r.extractWrap(h, opts)
	middlewares = append(middlewares, handler)
	r.app.Get(path, middlewares[0], middlewares[1:]...)
}

func (r *Router) Post(path string, h RouteHandler, opts ...RouterFuncOption) {
	handler, middlewares := r.extractWrap(h, opts)
	middlewares = append(middlewares, handler)
	r.app.Post(path, middlewares[0], middlewares[1:]...)
}

func (r *Router) Put(path string, h RouteHandler, opts ...RouterFuncOption) {
	handler, middlewares := r.extractWrap(h, opts)
	middlewares = append(middlewares, handler)
	r.app.Put(path, middlewares[0], middlewares[1:]...)
}

func (r *Router) Delete(path string, h RouteHandler, opts ...RouterFuncOption) {
	handler, middlewares := r.extractWrap(h, opts)
	middlewares = append(middlewares, handler)
	r.app.Delete(path, middlewares[0], middlewares[1:]...)
}

func (r *Router) Patch(path string, h RouteHandler, opts ...RouterFuncOption) {
	handler, middlewares := r.extractWrap(h, opts)
	middlewares = append(middlewares, handler)
	r.app.Patch(path, middlewares[0], middlewares[1:]...)
}

func (r *Router) Use(args ...any) IHttpRouter {
	r.app.Use(args...)
	return r
}

func (r *Router) Group(prefix string, handlers ...any) IHttpRouter {
	routeGroup := r.routeGroup
	var fiberHandlers []fiber.Handler

	if len(handlers) > 0 {
		for _, handler := range handlers {
			if name, ok := handler.(string); ok {
				routeGroup = name
			} else if fh, ok := handler.(fiber.Handler); ok {
				fiberHandlers = append(fiberHandlers, fh)
			} else if fh, ok := handler.(func(fiber.Ctx) error); ok {
				fiberHandlers = append(fiberHandlers, fh)
			}
		}
	}

	return &Router{
		app:                r.app.Group(prefix, fiberHandlers...),
		DisableAuthChecker: r.DisableAuthChecker,
		routeGroup:         routeGroup,
	}
}

// internal function to validate the access token
// this part is dynamic based on the usecase
func validateServiceToken(allowedService []string) error {
	if !slices.Contains(allowedService, config.GetConfig().GetString("service.name")) {
		return code.Get(code.UnauthorizedError).SetMessage("invalid service eligibility")
	}

	return nil
}

func validatePublicToken(routeRoles []string, allowedRoles []string) error {
	if len(routeRoles) == 0 {
		return nil
	}

	for _, role := range allowedRoles {
		if slices.Contains(routeRoles, role) {
			return nil
		}
	}

	return code.Get(code.UnauthorizedError).SetMessage("invalid role eligibility")
}
