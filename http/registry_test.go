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
				groupRegistry = map[string]RouteFunc{}
				RegisterRouteByGroup("test", []RouteFunc{
					func(r IHttpRouter) {
						r.Delete("/test", nil)
					},
				})
			},
			expectedLogs: []string{
				"DELETE /test/test",
			},
		},
		{
			name: "single public route",
			setup: func() {
				groupRegistry = map[string]RouteFunc{}

				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Get("/ping", nil) },
				})
			},
			expectedLogs: []string{
				"GET /ping",
			},
		},
		{
			name: "public and admin routes",
			setup: func() {
				groupRegistry = map[string]RouteFunc{} // reset
				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Get("/health", nil) },
				})
				RegisterRouteByGroup("admin", []RouteFunc{
					func(r IHttpRouter) { r.Post("/create", nil) },
				})
			},
			expectedLogs: []string{
				"GET /health",
				"POST /admin/create",
			},
		},
		{
			name: "multiple groups sorted",
			setup: func() {
				groupRegistry = map[string]RouteFunc{} // reset
				RegisterRouteByGroup("public", []RouteFunc{
					func(r IHttpRouter) { r.Delete("/test", nil) },
				})
				RegisterRouteByGroup("cms", []RouteFunc{
					func(r IHttpRouter) { r.Put("/test", nil) },
				})
				RegisterRouteByGroup("svc", []RouteFunc{
					func(r IHttpRouter) { r.Put("/test", nil) },
					func(r IHttpRouter) { r.Get("/test", nil) },
					func(r IHttpRouter) { r.Delete("/test", nil) },
					func(r IHttpRouter) { r.Post("/test", nil) },
					func(r IHttpRouter) { r.Patch("/test", nil) },
				})
			},
			expectedLogs: []string{
				"DELETE /test",
				"PUT /cms/test",
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
			mock := &MockRouter{Calls: &calls}
			tt.setup()
			RegisterRoutes(mock)

			got := append([]string(nil), *mock.Calls...)

			if !reflect.DeepEqual(got, tt.expectedLogs) {
				t.Errorf("unexpected router calls.\ngot:  %#v\nwant: %#v", mock.Calls, tt.expectedLogs)
			}
		})
	}
}
