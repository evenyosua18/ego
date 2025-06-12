package http

import (
	"sort"
)

var (
	groupRegistry = make(map[string]RouteFunc)
)

type (
	RouteFunc func(router IHttpRouter)
)

func RegisterRouteByGroup(group string, routeFns []RouteFunc) {
	registerGroup := func(group string, fn func(IHttpRouter)) {
		groupRegistry[group] = fn
	}

	registerGroup(group, func(r IHttpRouter) {
		for _, fn := range routeFns {
			fn(r)
		}
	})
}

func RegisterRoutes(r IHttpRouter) {
	// always register public routes first if available
	if fn, ok := groupRegistry["public"]; ok {
		fn(r.Group("/"))
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
		groupRegistry[groupName](r.Group("/" + groupName))
	}
}
