package router

var (
	groupRegistry = make(map[string]RouteFunc)
)

type (
	RouteFunc func(router IRouter)
)

func RegisterRouteByGroup(group string, routeFns []RouteFunc) {
	registerGroup := func(group string, fn func(IRouter)) {
		groupRegistry[group] = fn
	}

	registerGroup(group, func(r IRouter) {
		for _, fn := range routeFns {
			fn(r)
		}
	})
}

func RegisterRoutes(r IRouter) {
	for group, fn := range groupRegistry {
		if group == "public" {
			fn(r.Group("/"))
		} else {
			fn(r.Group("/" + group))
		}
	}
}
