package http

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFiberContext(t *testing.T) {
	type requestData struct {
		method string
		target string
		body   []byte
	}

	type expectedData struct {
		status     int
		jsonResult map[string]string
	}

	tests := []struct {
		name     string
		request  requestData
		expected expectedData
		setup    func(t *testing.T, app *fiber.App)
	}{
		{
			name: "test all context functions",
			request: requestData{
				method: http.MethodPost,
				target: "/test/123?name=ego",
				body:   []byte(`{"message":"hello"}`),
			},
			expected: expectedData{
				status: 200,
				jsonResult: map[string]string{
					"message": "success",
				},
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Post("/test/:id", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}

					// Param
					if got := ctx.Param("id"); got != "123" {
						t.Errorf("Param() = %v, want %v", got, "123")
					}

					// Query
					if got := ctx.Query("name"); got != "ego" {
						t.Errorf("Query() = %v, want %v", got, "ego")
					}

					// Body
					expectedBody := []byte(`{"message":"hello"}`)
					if got := ctx.Body(); !bytes.Equal(got, expectedBody) {
						t.Errorf("Body() = %s, want %s", string(got), string(expectedBody))
					}

					// Respond with JSON
					return ctx.JSON(200, map[string]string{"message": "success"})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			// Setup routes
			tt.setup(t, app)

			req := httptest.NewRequest(tt.request.method, tt.request.target, bytes.NewReader(tt.request.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expected.status {
				t.Errorf("StatusCode = %v, want %v", resp.StatusCode, tt.expected.status)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var actual map[string]string
			if err := json.Unmarshal(respBody, &actual); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			for key, want := range tt.expected.jsonResult {
				if got := actual[key]; got != want {
					t.Errorf("JSON response[%s] = %v, want %v", key, got, want)
				}
			}
		})
	}
}
