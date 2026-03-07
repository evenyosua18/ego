package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/request"
)

// ValidateToken calls the authorization service to validate the given access token
func ValidateToken(ctx context.Context, accessToken string) (*ClaimToken, error) {
	client := request.NewClient(nil)

	reqBody := map[string]string{
		"access_token": accessToken,
	}

	url := config.GetConfig().GetString("auth_svc.base_url") + "/token/validate"

	req := request.Request{
		Method: http.MethodPost,
		URL:    url,
		Body:   reqBody,
	}

	resp, body, err := client.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call validate token api: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("validate token api returned status: %d", resp.StatusCode)
	}

	var claim ClaimToken
	if err := json.Unmarshal(body, &claim); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claim token: %w", err)
	}

	return &claim, nil
}
