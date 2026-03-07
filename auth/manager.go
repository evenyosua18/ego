package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/evenyosua18/ego/logger"
	"github.com/evenyosua18/ego/request"
)

var (
	accessToken   string
	tokenMutex    sync.RWMutex
	managerCtx    context.Context
	managerCancel context.CancelFunc
)

type generateTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type generateTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GetAccessToken securely returns the active access token.
func GetAccessToken() string {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return accessToken
}

func setAccessToken(token string) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	accessToken = token
}

// ManageAccessToken kicks off the background token manager.
// It fetches an initial token synchronously so the service is ready immediately,
// and starts a background goroutine to refresh it before expiration.
func ManageAccessToken(ctx context.Context, baseUrl, clientID, clientSecret string) error {
	// Setup context for the background routine
	managerCtx, managerCancel = context.WithCancel(ctx)

	// Fetch the initial token synchronously
	expiredAt, err := fetchAndSetToken(managerCtx, baseUrl, clientID, clientSecret)
	if err != nil {
		managerCancel()
		return fmt.Errorf("failed fetching initial access token: %w", err)
	}

	logger.Info(fmt.Sprintf("Auth Manager: Initial access token fetched successfully. Expires at %v", expiredAt))

	// Start background refresh routine
	go refreshRoutine(baseUrl, clientID, clientSecret, expiredAt)

	return nil
}

func fetchAndSetToken(ctx context.Context, baseUrl, clientID, clientSecret string) (time.Time, error) {
	genUrl := baseUrl + "v1/token/generate"
	reqBody := generateTokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	client := request.NewClient(nil)
	req := request.Request{
		Method: http.MethodPost,
		URL:    genUrl,
		Body:   reqBody,
	}

	resp, body, err := client.Do(ctx, req)
	if err != nil {
		return time.Time{}, fmt.Errorf("generate token request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return time.Time{}, fmt.Errorf("generate token returned status: %d", resp.StatusCode)
	}

	var genResp generateTokenResponse
	if err := json.Unmarshal(body, &genResp); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode generate token response: %w", err)
	}

	// Now validate to get the expiration time
	validateUrl := baseUrl + "v1/token/validate"
	valReqBody := map[string]string{
		"access_token": genResp.AccessToken,
	}
	valReq := request.Request{
		Method: http.MethodPost,
		URL:    validateUrl,
		Body:   valReqBody,
	}

	valResp, valBody, err := client.Do(ctx, valReq)
	if err != nil {
		return time.Time{}, fmt.Errorf("validate token request failed: %w", err)
	}

	if valResp.StatusCode < 200 || valResp.StatusCode >= 300 {
		return time.Time{}, fmt.Errorf("validate token returned status: %d", valResp.StatusCode)
	}

	var claim ClaimToken
	if err := json.Unmarshal(valBody, &claim); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode validate token response: %w", err)
	}

	// Parse the expiration time
	expiredAt, err := time.Parse(time.RFC3339, claim.ExpiredAt)
	if err != nil {
		// Fallback if parsing fails
		logger.Error(fmt.Errorf("failed to parse expiration time %s: %v", claim.ExpiredAt, err))
		expiredAt = time.Now().Add(50 * time.Minute)
	}

	// Securely set the new token
	setAccessToken(genResp.AccessToken)

	return expiredAt, nil
}

func refreshRoutine(baseUrl, clientID, clientSecret string, initialExpiredAt time.Time) {
	expiredAt := initialExpiredAt

	for {
		// Calculate time until refresh (e.g., 5 minutes before it actually expires)
		// If the token is already expired or close to it, refresh immediately
		refreshTime := expiredAt.Add(-5 * time.Minute)
		durationUntilRefresh := time.Until(refreshTime)

		if durationUntilRefresh <= 0 {
			durationUntilRefresh = 5 * time.Second // Retry shortly if we are already past the refresh mark
		}

		select {
		case <-managerCtx.Done():
			logger.Info("Auth Manager: Shutting down background token refresh.")
			return
		case <-time.After(durationUntilRefresh):
			logger.Info("Auth Manager: Refreshing access token...")
			newExpiredAt, err := fetchAndSetToken(managerCtx, baseUrl, clientID, clientSecret)
			if err != nil {
				logger.Error(fmt.Errorf("Auth Manager: Failed to refresh token: %v. Will retry in 1 minute.", err))
				// If refresh fails, retry in 1 minute instead of waiting for full expiration again
				expiredAt = time.Now().Add(6 * time.Minute)
			} else {
				logger.Info(fmt.Sprintf("Auth Manager: Token refreshed successfully. Next expiration at %v", newExpiredAt))
				expiredAt = newExpiredAt
			}
		}
	}
}
