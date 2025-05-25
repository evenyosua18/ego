package http

type (
	RouterOptions struct {
		Roles []string
	}

	RouterFuncOption func(opt *RouterOptions)
)

func SetRouterFullOption(routerOpts RouterOptions) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		opt.Roles = append(opt.Roles, routerOpts.Roles...) // IMPROVEMENT: remove duplicates and normalize

		// NOTE: add here if there are new router options
	})
}

func SetRouterRolesOption(roles []string) RouterFuncOption {
	return RouterFuncOption(func(opt *RouterOptions) {
		opt.Roles = append(opt.Roles, roles...) // IMPROVEMENT: remove duplicates and normalize
	})
}
