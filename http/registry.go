package http

import (
	"sort"
)

type (
	RouteFunc func(router IHttpRouter)
	RouteGroup struct {
		Name     string
		Handlers []any
		RouteFns []RouteFunc
	}
)

var (
	groupRegistry = make(map[string]RouteGroup)
)

func RegisterRouteByGroup(group string, routeFns []RouteFunc, handlers ...any) {
	registerGroup := func(group string, fn []RouteFunc, handlers []any) {
		groupRegistry[group] = RouteGroup{
			Name:     group,
			Handlers: handlers,
			RouteFns: fn,
		}
	}

	registerGroup(group, routeFns, handlers)
}

func RegisterRoutes(r IHttpRouter) {
	// always register public routes first if available
	if g, ok := groupRegistry["public"]; ok {
		// prepend group name to handlers
		handlers := append([]any{"public"}, g.Handlers...)
		routerGroup := r.Group("/", handlers...)
		for _, fn := range g.RouteFns {
			fn(routerGroup)
		}
	}

	// map groups, except public
	var groups []string
	for g := range groupRegistry {
		if g != "public" {
			groups = append(groups, g)
		}
	}

	// sort group
	sort.Strings(groups)

	// loop based on sorted group name
	for _, groupName := range groups {
		g := groupRegistry[groupName]
		handlers := append([]any{groupName}, g.Handlers...)
		routerGroup := r.Group("/"+groupName, handlers...)
		for _, fn := range g.RouteFns {
			fn(routerGroup)
		}
	}
}
