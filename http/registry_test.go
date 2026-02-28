package http

import (
	"reflect"
	"testing"
)

func TestRegisterRoutes(t *testing.T) {
	tests := []struct {
		name         string
		setup        func()
		expectedLogs []string
	}{
		{
			name: "route under group",
			setup: func() {
				groupRegistry = map[string]RouteGroup{}
				RegisterRouteByGroup("test", []RouteFunc{
					func(r IHttpRouter) {
						r.Delete("/test", nil)
					},
				})
			},
			expectedLogs: []string{
				"GROUP /test AS test",
				"DELETE /test/test",
			},
		},
		{
			name: "single public route",
			setup: func() {
				groupRegistry = map[string]RouteGroup{}

				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Get("/ping", nil) },
				})
			},
			expectedLogs: []string{
				"GROUP / AS public",
				"GET /ping",
			},
		},
		{
			name: "public and admin routes",
			setup: func() {
				groupRegistry = map[string]RouteGroup{} // reset
				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Get("/health", nil) },
				})
				RegisterRouteByGroup("admin", []RouteFunc{
					func(r IHttpRouter) { r.Post("/create", nil) },
				})
			},
			expectedLogs: []string{
				"GROUP / AS public",
				"GET /health",
				"GROUP /admin AS admin",
				"POST /admin/create",
			},
		},
		{
			name: "multiple groups sorted",
			setup: func() {
				groupRegistry = map[string]RouteGroup{} // reset
				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Delete("/test", nil) },
				})
				RegisterRouteByGroup("cms", []RouteFunc{
					func(r IHttpRouter) { r.Put("/test", nil) },
				}, "middleware1", "middleware2")
				RegisterRouteByGroup("svc", []RouteFunc{
					func(r IHttpRouter) { r.Put("/test", nil) },
					func(r IHttpRouter) { r.Get("/test", nil) },
					func(r IHttpRouter) { r.Delete("/test", nil) },
					func(r IHttpRouter) { r.Post("/test", nil) },
					func(r IHttpRouter) { r.Patch("/test", nil) },
				})
			},
			expectedLogs: []string{
				"GROUP / AS public",
				"DELETE /test",
				"GROUP /cms AS cms",
				"GROUP /cms AS middleware1",
				"GROUP /cms AS middleware2",
				"PUT /cms/test",
				"GROUP /svc AS svc",
				"PUT /svc/test",
				"GET /svc/test",
				"DELETE /svc/test",
				"POST /svc/test",
				"PATCH /svc/test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls []string
			mockHttpRouter := &HttpRouter{Calls: &calls}
			tt.setup()
			RegisterRoutes(mockHttpRouter)

			got := append([]string(nil), *mockHttpRouter.Calls...)

			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("unexpected router calls.\ngot:  %#v\nwant: %#v", mockHttpRouter.Calls, tt.expectedLogs)
			}
		})
	}
}
