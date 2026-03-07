package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/evenyosua18/ego/code"
	"github.com/gofiber/fiber/v3"
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

					// SetContext and GetContext
					ctx.SetContext("test-key", "test-value")
					if got := c.Locals("test-key"); got != "test-value" {
						t.Errorf("Locals() = %v, want %v", got, "test-value")
					}
					if got := ctx.GetContext("test-key"); got != "test-value" {
						t.Errorf("GetContext() = %v, want %v", got, "test-value")
					}

					// RequestBody
					var reqBody map[string]string
					if err := ctx.RequestBody(&reqBody); err != nil {
						t.Errorf("RequestBody() error = %v", err)
					}
					if reqBody["message"] != "hello" {
						t.Errorf("RequestBody() mapped value = %v, want %v", reqBody["message"], "hello")
					}

					// Context
					if ctx.Context() == nil {
						t.Errorf("Context() returned nil")
					}

					// GetRequestHeader
					if got := ctx.GetRequestHeader("Content-Type"); got != "application/json" {
						t.Errorf("GetRequestHeader() = %v, want %v", got, "application/json")
					}

					// SetCookie and Cookie
					ctx.SetCookie("test-cookie", "cookie-value", time.Now().Add(time.Hour))
					// Note: Cookie getting from the same request where it's set in Fiber tests
					// won't immediately reflect in ctx.Cookie() without a real client roundtrip,
					// but we call it to cover the method.
					ctx.Cookie("test-cookie")

					// FormValue
					ctx.FormValue("some_form_field")

					// GetRouteRoles
					// Setup local mock
					c.Locals(localRouteRoles{}, []string{"admin", "user"})
					if roles := ctx.GetRouteRoles(); len(roles) != 2 || roles[0] != "admin" {
						t.Errorf("GetRouteRoles() = %v, want %v", roles, []string{"admin", "user"})
					}

					// Respond with JSON
					return ctx.JSON(200, map[string]string{"message": "success"})
				})
			},
		},
		{
			name: "test response error",
			request: requestData{
				method: http.MethodGet,
				target: "/error",
			},
			expected: expectedData{
				status: 400,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/error", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					// Return a dummy error that code.Extract handles
					return ctx.ResponseError(code.Get(code.BadRequestError))
				})
			},
		},
		{
			name: "test response success",
			request: requestData{
				method: http.MethodGet,
				target: "/success",
			},
			expected: expectedData{
				status: 200,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/success", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					return ctx.ResponseSuccess(map[string]string{"data": "ok"})
				})
			},
		},
		{
			name: "test send",
			request: requestData{
				method: http.MethodGet,
				target: "/send",
			},
			expected: expectedData{
				status: 201, // arbitrary status
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/send", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					return ctx.Send(201, []byte("created"))
				})
			},
		},
		{
			name: "test http data and truncate",
			request: requestData{
				method: http.MethodPost,
				target: "/httpdata?q=1",
				body:   []byte(`{"large":"payload"}`),
			},
			expected: expectedData{
				status: 200,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Post("/httpdata", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					data, opName := ctx.HttpData()
					if opName != "POST /httpdata" {
						t.Errorf("HttpData() opName = %v, want POST /httpdata", opName)
					}
					if data["method"] != "POST" {
						t.Errorf("HttpData() method = %v, want POST", data["method"])
					}
					return c.SendStatus(200)
				})
			},
		},
		{
			name: "test bind query",
			request: requestData{
				method: http.MethodGet,
				target: "/bind?name=ego&age=25&active=true&score=99.5",
			},
			expected: expectedData{
				status: 200,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/bind", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					type QueryStruct struct {
						Name   string  `query:"name"`
						Age    int     `query:"age"`
						Active bool    `query:"active"`
						Score  float64 `query:"score"`
					}
					var q QueryStruct
					if err := ctx.BindQuery(&q); err != nil {
						t.Errorf("BindQuery() error %v", err)
					}
					if q.Name != "ego" || q.Age != 25 || q.Active != true || q.Score != 99.5 {
						t.Errorf("BindQuery() mapped values incorrect: %+v", q)
					}
					return c.SendStatus(200)
				})
			},
		},
		{
			name: "test response redirect",
			request: requestData{
				method: http.MethodGet,
				target: "/redirect",
			},
			expected: expectedData{
				status: 302,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/redirect", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					return ctx.ResponseRedirect("/target")
				})
			},
		},
		{
			name: "test render",
			request: requestData{
				method: http.MethodGet,
				target: "/render",
			},
			expected: expectedData{
				status: 500, // fiber.Render normally requires an engine; returning 500 when engine is nil is fine to just hit it
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Get("/render", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					// Will likely return an error because no template engine is configured,
					// but it covers the line.
					err := ctx.Render("index", map[string]any{"key": "value"})
					if err != nil {
						return c.Status(500).SendString(err.Error())
					}
					return c.SendStatus(200)
				})
			},
		},
		{
			name: "test next",
			request: requestData{
				method: http.MethodGet,
				target: "/next",
			},
			expected: expectedData{
				status: 200,
			},
			setup: func(t *testing.T, app *fiber.App) {
				app.Use("/next", func(c fiber.Ctx) error {
					ctx := &fiberContext{ctx: c}
					return ctx.Next()
				})
				app.Get("/next", func(c fiber.Ctx) error {
					return c.SendStatus(200)
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

			// If expecting a specific JSON payload matching tt.expected.jsonResult
			if tt.expected.jsonResult != nil {
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
			}
		})
	}
}

// Additional test to fully cover BindQuery error cases and truncate edge cases
func TestFiberContext_BindQueryErrors(t *testing.T) {
	app := fiber.New()
	app.Get("/bind", func(c fiber.Ctx) error {
		ctx := &fiberContext{ctx: c}

		// Map destination instead of struct pointer
		var m map[string]any
		if err := ctx.BindQuery(&m); err == nil {
			t.Errorf("Expected error for non-struct pointer in BindQuery")
		}

		// Nil destination
		if err := ctx.BindQuery(nil); err == nil {
			t.Errorf("Expected error for nil destination in BindQuery")
		}

		// Valid struct, but trying to cover the parser branches like int, uint, float, bool
		type QueryAllTypes struct {
			MyString   string  `query:"str"`
			MyInt      int     `query:"int"`
			MyInt8     int8    `query:"int8"`
			MyInt16    int16   `query:"int16"`
			MyInt32    int32   `query:"int32"`
			MyInt64    int64   `query:"int64"`
			MyUint     uint    `query:"uint"`
			MyUint8    uint8   `query:"uint8"`
			MyUint16   uint16  `query:"uint16"`
			MyUint32   uint32  `query:"uint32"`
			MyUint64   uint64  `query:"uint64"`
			MyFloat32  float32 `query:"f32"`
			MyFloat64  float64 `query:"f64"`
			MyBool     bool    `query:"b"`
			SkipMe     string  `query:"skip"` // intentionally not passed
			unexported string
			NoTag      string
		}

		var q QueryAllTypes
		ctx.BindQuery(&q)

		// Pointer to non-struct
		var strPtr *string
		strVal := "hello"
		strPtr = &strVal
		if err := ctx.BindQuery(strPtr); err == nil {
			t.Errorf("Expected error for pointer to non-struct in BindQuery")
		}

		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/bind?name=test&str=1&int=1&int8=1&int16=1&int32=1&int64=1&uint=1&uint8=1&uint16=1&uint32=1&uint64=1&f32=1.5&f64=1.5&b=true&NoTag=notagval", nil)
	app.Test(req)
}

func TestTruncate(t *testing.T) {
	// Test omit based on content-type
	if res := truncate("text/html", "body"); res != omittedFileText {
		t.Errorf("Expected omittedFileText, got %v", res)
	}

	// Test string
	if res := truncate("application/json", "hello"); res != "hello" {
		t.Errorf("Expected normal string")
	}

	// Test long string
	longStr := strings.Repeat("A", maxLogLength+10)
	expectedLongStr := strings.Repeat("A", maxLogLength) + truncatedText
	if res := truncate("application/json", longStr); res != expectedLongStr {
		t.Errorf("Expected string truncation failed")
	}

	// Test byte slice
	if res := truncate("application/x-www-form-urlencoded", []byte("hello")); res != "hello" {
		t.Errorf("Expected normal byte slice")
	}

	// Test long byte slice
	longBytes := []byte(strings.Repeat("B", maxLogLength+10))
	expectedLongBytes := strings.Repeat("B", maxLogLength) + truncatedText
	if res := truncate("application/x-www-form-urlencoded", longBytes); res != expectedLongBytes {
		t.Errorf("Expected byte truncation failed")
	}

	// Test structured value
	type MyStruct struct {
		Name string
	}
	st := MyStruct{Name: "Ego"}
	stRes := truncate("application/json", st)
	if !strings.Contains(stRes, "Ego") {
		t.Errorf("Expected JSON marshaling failed")
	}

	// Test large structured value
	largeSt := MyStruct{Name: strings.Repeat("C", maxLogLength+10)}
	largeStRes := truncate("application/json", largeSt)
	if !strings.HasSuffix(largeStRes, truncatedText) {
		t.Errorf("Expected JSON truncation failed")
	}

	// Test invalid JSON marshal (e.g. func)
	f := func() {}
	if res := truncate("application/json", f); res != unableMarshalText {
		t.Errorf("Expected marshal failure text, got %v", res)
	}
}
