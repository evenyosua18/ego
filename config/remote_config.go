package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type RemoteConfigProvider interface {
	Fetch(ctx context.Context) (map[string]any, error)
}

// FirebaseRemoteConfig implements RemoteConfigProvider for a generic Firebase URL returning JSON
type FirebaseRemoteConfig struct{
	credentials []byte
}

func NewFirebaseRemoteConfig(credentials []byte) *FirebaseRemoteConfig {
	return &FirebaseRemoteConfig{
		credentials: credentials,
	}
}

func (f *FirebaseRemoteConfig) Fetch(ctx context.Context) (map[string]any, error) {
	var opts []option.ClientOption
	if len(f.credentials) > 0 {
		opts = append(opts, option.WithCredentialsJSON(f.credentials))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.RemoteConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize remote config client: %w", err)
	}

	// Fetch server template natively from Firebase Admin SDK
	template, err := client.GetServerTemplate(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server template: %w", err)
	}

	templateJSON, err := template.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize server template: %w", err)
	}

	var parsedTemplate struct {
		Parameters map[string]struct {
			DefaultValue struct {
				Value *string `json:"value"`
			} `json:"defaultValue"`
		} `json:"parameters"`
	}

	if err := json.Unmarshal([]byte(templateJSON), &parsedTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse template json: %w", err)
	}

	result := make(map[string]any)
	for key, param := range parsedTemplate.Parameters {
		if param.DefaultValue.Value != nil {
			valStr := *param.DefaultValue.Value
			var valAny any
			if err := json.Unmarshal([]byte(valStr), &valAny); err != nil {
				valAny = valStr
			}
			result[key] = valAny
		}
	}

	return result, nil
}

// AutoRefresh starts a background goroutine that fetches and updates configuration periodically.
func AutoRefresh(ctx context.Context, provider RemoteConfigProvider, period time.Duration, updateFn func(map[string]any)) {
	if period <= 0 {
		return
	}
	
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				values, err := provider.Fetch(ctx)
				if err != nil {
					// We just skip this refresh cycle on error.
					// In a real scenario, we might want to log this using ego logger,
					// but here we depend on the updateFn safely updating the config.
					continue
				}
				if updateFn != nil {
					updateFn(values)
				}
			}
		}
	}()
}
