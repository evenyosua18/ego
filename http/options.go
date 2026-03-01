package http

import fiber "github.com/gofiber/fiber/v3"

type (
	RouterOptions struct {
		Roles       []string
		Middlewares []fiber.Handler
	}

	RouterFuncOption func(opt *RouterOptions)
)

func SetRouterFullOption(routerOpts RouterOptions) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		opt.Roles = append(opt.Roles, routerOpts.Roles...) // IMPROVEMENT: remove duplicates and normalize
		opt.Middlewares = append(opt.Middlewares, routerOpts.Middlewares...)

		// NOTE: add here if there are new router options
	})
}

func SetRouterRolesOption(roles []string) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		opt.Roles = append(opt.Roles, roles...) // IMPROVEMENT: remove duplicates and normalize
	})
}

func SetRouterMiddleware(middlewares ...fiber.Handler) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		opt.Middlewares = append(opt.Middlewares, middlewares...)
	})
}

func SetMiddleware(middlewares ...RouteHandler) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		for _, m := range middlewares {
			handler := m
			opt.Middlewares = append(opt.Middlewares, func(c fiber.Ctx) error {
				return handler(&fiberContext{ctx: c})
			})
		}
	})
}
