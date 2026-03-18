package request

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Methods_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Output headers and method back to verify
		w.Header().Set("X-Received-Method", r.Method)
		w.Header().Set("X-Received-Query", r.URL.Query().Get("user"))
		w.Header().Set("X-Received-Custom", r.Header.Get("X-Custom-Auth"))

		bodyBytes, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bodyBytes) // Echo body
	}))
	defer ts.Close()

	client := NewClient(ts.Client())

	opts := []RequestOption{
		WithHeaders(map[string]string{"X-Custom-Auth": "secret"}),
		WithQueryParams(map[string]string{"user": "123"}),
	}

	tests := []struct {
		name       string
		methodFunc func(context.Context, string, ...RequestOption) (*http.Response, []byte, error)
		reqMethod  string
		reqBody    any
	}{
		{
			name:       "Get Method",
			methodFunc: client.Get,
			reqMethod:  http.MethodGet,
			reqBody:    nil,
		},
		{
			name:       "Post Method",
			methodFunc: client.Post,
			reqMethod:  http.MethodPost,
			reqBody:    map[string]string{"msg": "hello post"},
		},
		{
			name:       "Put Method",
			methodFunc: client.Put,
			reqMethod:  http.MethodPut,
			reqBody:    []byte(`{"msg":"hello put"}`),
		},
		{
			name:       "Delete Method",
			methodFunc: client.Delete,
			reqMethod:  http.MethodDelete,
			reqBody:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOpts := opts
			if tt.reqBody != nil {
				currentOpts = append(currentOpts, WithBody(tt.reqBody))
			}

			resp, body, err := tt.methodFunc(context.Background(), ts.URL, currentOpts...)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if resp.Header.Get("X-Received-Method") != tt.reqMethod {
				t.Errorf("Expected method %s, got %s", tt.reqMethod, resp.Header.Get("X-Received-Method"))
			}

			if resp.Header.Get("X-Received-Query") != "123" {
				t.Errorf("Expected query parameter user=123")
			}

			if resp.Header.Get("X-Received-Custom") != "secret" {
				t.Errorf("Expected custom header secret")
			}

			if tt.reqBody != nil {
				var expectedBody string
				if b, ok := tt.reqBody.([]byte); ok {
					expectedBody = string(b)
				} else {
					jsonBytes, _ := json.Marshal(tt.reqBody)
					expectedBody = string(jsonBytes)
				}

				if string(body) != expectedBody {
					t.Errorf("Expected body echo %s, got %s", expectedBody, string(body))
				}
			}
		})
	}
}

func TestClient_Retry_SuccessAfterFails(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError) // 500 should trigger retry
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer ts.Close()

	client := NewClient(ts.Client())

	// Request should attempt 3 times: fail, fail, success
	resp, body, err := client.Get(context.Background(), ts.URL, WithRetry(2, 5*time.Millisecond, http.StatusInternalServerError))
	
	if err != nil {
		t.Fatalf("Unexpected error after retries: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected final status OK, got %d", resp.StatusCode)
	}

	if string(body) != "success" {
		t.Errorf("Expected body success, got %s", string(body))
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestClient_Retry_Exhausted(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusTooManyRequests) // 429 should trigger retry
	}))
	defer ts.Close()

	client := NewClient(ts.Client())

	// Request should attempt 4 times: fail 4 times
	resp, _, err := client.Get(context.Background(), ts.URL, WithRetry(3, 5*time.Millisecond))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected final status 429, got %d", resp.StatusCode)
	}

	if attempts != 4 {
		t.Errorf("Expected 4 attempts, got %d", attempts)
	}
}

func TestClient_Retry_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest) // 400 should NOT trigger retry
	}))
	defer ts.Close()

	client := NewClient(ts.Client())

	// Request should attempt only 1 time
	resp, _, err := client.Get(context.Background(), ts.URL, WithRetry(3, 5*time.Millisecond))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected final status 400, got %d", resp.StatusCode)
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestClient_Do_Success_GET(t *testing.T) {
	// Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Query().Get("key") != "value" {
			t.Errorf("Expected query parameter key=value, got %s", r.URL.Query().Get("key"))
		}
		if r.Header.Get("Custom-Header") != "TestValue" {
			t.Errorf("Expected custom header, got %s", r.Header.Get("Custom-Header"))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success"}`))
	}))
	defer ts.Close()

	client := NewClient(ts.Client())
	req := Request{
		Method: http.MethodGet,
		URL:    ts.URL,
		Headers: map[string]string{
			"Custom-Header": "TestValue",
		},
		QueryParams: map[string]string{
			"key": "value",
		},
	}

	resp, respBody, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}

	expectedBody := `{"status":"success"}`
	if string(respBody) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(respBody))
	}
}

func TestClient_Do_Success_POST_JSON(t *testing.T) {
	type BodyPayload struct {
		Name string `json:"name"`
	}

	// Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}
		
		var payload BodyPayload
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Fatalf("Failed to unmarshal body: %v", err)
		}

		if payload.Name != "TestName" {
			t.Errorf("Expected body name TestName, got %s", payload.Name)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"created":true}`))
	}))
	defer ts.Close()

	client := NewClient(ts.Client())
	req := Request{
		Method: http.MethodPost,
		URL:    ts.URL,
		Body: BodyPayload{
			Name: "TestName",
		},
	}

	resp, respBody, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %d", resp.StatusCode)
	}

	expectedBody := `{"created":true}`
	if string(respBody) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(respBody))
	}
}

func TestClient_Do_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate slow response
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ts.Client())
	req := Request{
		Method: http.MethodGet,
		URL:    ts.URL,
	}

	// Setup context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, _, err := client.Do(ctx, req)
	if err == nil {
		t.Fatal("Expected error due to context timeout, but got nil")
	}
}

func TestClient_Do_ErrorStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad_request"}`))
	}))
	defer ts.Close()

	client := NewClient(ts.Client())
	req := Request{
		Method: http.MethodGet,
		URL:    ts.URL,
	}

	resp, respBody, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Did not expect an error from client.Do() for 400 status codes, got: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	expectedBody := `{"error":"bad_request"}`
	if string(respBody) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(respBody))
	}
}

func TestClient_CircuitBreakerAndCheckEndpoints(t *testing.T) {
	hitCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hitCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient(ts.Client())

	// The breaker trips after > 5 consecutive failures (so 6 failures).
	// We will make 6 requests that fail.
	for i := 0; i < 6; i++ {
		_, _, _ = client.Get(context.Background(), ts.URL)
	}

	// At this point, the circuit breaker should be Open.
	// 1. Verify CheckEndpoints returns an error.
	err := CheckEndpoints(ts.URL)
	if err == nil {
		t.Error("Expected CheckEndpoints to return an error (circuit should be open)")
	}

	// 2. Verify subsequent requests fail instantly without hitting the server.
	hitBefore := hitCount
	_, _, err = client.Get(context.Background(), ts.URL)
	if err == nil {
		t.Error("Expected request to return error when circuit is open")
	}

	if hitCount > hitBefore {
		t.Errorf("Expected no additional hits to server, but got %d extra hits", hitCount-hitBefore)
	}
}
