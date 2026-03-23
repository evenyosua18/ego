package config

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRemoteProvider struct {
	values map[string]any
	err    error
	calls  int
}

func (m *mockRemoteProvider) Fetch(ctx context.Context) (map[string]any, error) {
	m.calls++
	return m.values, m.err
}

func TestAutoRefresh(t *testing.T) {
	provider := &mockRemoteProvider{
		values: map[string]any{"test_key": "test_val"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updatedValues := make(map[string]any)
	updateFn := func(v map[string]any) {
		for key, val := range v {
			updatedValues[key] = val
		}
	}

	AutoRefresh(ctx, provider, 10*time.Millisecond, updateFn)

	time.Sleep(35 * time.Millisecond)

	assert.GreaterOrEqual(t, provider.calls, 2)
	assert.Equal(t, "test_val", updatedValues["test_key"])
}
