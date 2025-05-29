package http

import (
	"github.com/gofiber/fiber/v3"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Methods(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		registerFn func(r *Router)
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:   "GET method",
			method: http.MethodGet,
			registerFn: func(r *Router) {
				r.Get("/ping", func(c Context) error {
					return c.Send(200, []byte("pong"))
				})
			},
			path:       "/ping",
			wantStatus: http.StatusOK,
			wantBody:   "pong",
		},
		{
			name:   "POST method",
			method: http.MethodPost,
			registerFn: func(r *Router) {
				r.Post("/post", func(c Context) error {
					return c.Send(200, []byte("TEST"))
				})
			},
			path:       "/post",
			wantStatus: http.StatusOK,
			wantBody:   "TEST",
		},
		{
			name:   "PUT method",
			method: http.MethodPut,
			registerFn: func(r *Router) {
				r.Put("/put", func(c Context) error {
					return c.Send(200, []byte("TEST"))
				})
			},
			path:       "/put",
			wantStatus: http.StatusOK,
			wantBody:   "TEST",
		},
		{
			name:   "DELETE method",
			method: http.MethodDelete,
			registerFn: func(r *Router) {
				r.Delete("/delete", func(c Context) error {
					return c.Send(200, []byte("TEST"))
				})
			},
			path:       "/delete",
			wantStatus: http.StatusOK,
			wantBody:   "TEST",
		},
		{
			name:   "PATCH method",
			method: http.MethodPatch,
			registerFn: func(r *Router) {
				r.Patch("/patch", func(c Context) error {
					return c.Send(200, []byte("TEST"))
				})
			},
			path:       "/patch",
			wantStatus: http.StatusOK,
			wantBody:   "TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup router
			app := fiber.New()
			router := &Router{app: app}
			tt.registerFn(router)

			// setup new request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status code = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			// Check response body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			if got := string(respBody); got != tt.wantBody {
				t.Errorf("response body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}

func TestRouter_Group(t *testing.T) {
	// setup router
	app := fiber.New()
	r := &Router{app: app}

	group := r.Group("/api")
	group.Get("/hello", func(c Context) error {
		return c.JSON(200, fiber.Map{"msg": "hello from group"})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

}
