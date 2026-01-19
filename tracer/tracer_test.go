package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracerSingleton(t *testing.T) {
	// 1. Get a reference to the tracer BEFORE Config
	t1 := GetTracer()

	// Verify t1 is initially inactive (default state)
	sentryTracer1, ok := t1.(*SentryTracer)
	assert.True(t, ok)
	assert.False(t, sentryTracer1.isActive)

	// 2. RunSentry with valid config
	config := Config{
		Dsn:             "https://examplePublicKey@o0.ingest.sentry.io/0",
		Env:             "test",
		TraceSampleRate: 1.0,
	}
	// We don't actually want to init real sentry network calls if possible,
	// but RunSentry calls sentry.Init.
	// For this test, we care about the variable reference.
	// sentry.Init might fail if DSN is invalid, but we can try mocking or ignoring error if it's just about variable assignment.
	// However, RunSentry returns error if Init fails.

	// To avoid actual network, we might receive error from sentry.Init or it might pass validation.
	// Let's assume validation passes for format.

	_, _ = RunSentry(config) // Ignore error, we just want to see if variable mutated.

	// 3. Get reference again
	t2 := GetTracer()

	// 4. Verify pointers are the same
	assert.Equal(t, t1, t2, "Tracer references should be identical")

	// 5. Verify t1 is now active (mutation happened)
	assert.True(t, sentryTracer1.isActive, "Tracer should be active after RunSentry")
}
