package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestManageAccessToken(t *testing.T) {
	// 1. Create a mock auth_svc
	generateCount := 0
	validateCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/token/generate" {
			generateCount++
			resp := map[string]string{
				"access_token":  "mock-token-123",
				"refresh_token": "mock-refresh-123",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		if r.URL.Path == "/v1/token/validate" {
			validateCount++
			// Set expiration to 6 seconds from now so the 5 minute refresh window triggers quickly
			expireTime := time.Now().Add(6 * time.Second)
			resp := map[string]interface{}{
				"expired_at": expireTime.Format(time.RFC3339),
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	// 2. Call ManageAccessToken
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := ManageAccessToken(ctx, ts.URL+"/", "client123", "secret123")
	if err != nil {
		t.Fatalf("Failed to initialize token manager: %v", err)
	}

	// 3. Verify synchronous fetch succeeded immediately
	token := GetAccessToken()
	if token != "mock-token-123" {
		t.Errorf("Expected token mock-token-123, got %s", token)
	}
	if generateCount != 1 {
		t.Errorf("Expected 1 call to /generate, got %d", generateCount)
	}
	if validateCount != 1 {
		t.Errorf("Expected 1 call to /validate, got %d", validateCount)
	}

	// 4. Wait a couple of seconds and see if refresh happens
	// The original token lives for 6 seconds, we refresh 5 minutes before.
	// Since 6 seconds is < 5 minutes, it should retry very quickly (5 second delay per our code).
	time.Sleep(6 * time.Second)

	if generateCount < 2 {
		t.Errorf("Expected token to have been refreshed automatically. generateCount: %d", generateCount)
	}
}
