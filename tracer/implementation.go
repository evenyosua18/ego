package tracer

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"time"
)

var (
	errConfigValue = fmt.Errorf("configuration value is invalid")

	tracer Tracer = &SentryTracer{isActive: false}
)

type (
	Config struct {
		Dsn             string
		Env             string
		FlushTime       string
		TraceSampleRate float64
	}
)

func RunSentry(sentryConfig Config) (flush func(flushTime string), err error) {
	if sentryConfig.Dsn == "" || sentryConfig.Env == "" {
		return nil, errConfigValue
	}

	if sentryConfig.TraceSampleRate > 1.0 {
		sentryConfig.TraceSampleRate = 1.0
	}

	if err = sentry.Init(sentry.ClientOptions{
		Dsn:              sentryConfig.Dsn,
		Environment:      sentryConfig.Env,
		TracesSampleRate: sentryConfig.TraceSampleRate,
		Transport:        sentry.NewHTTPTransport(),
		EnableTracing:    true,
	}); err != nil {
		return nil, err
	}

	tracer = &SentryTracer{isActive: true}

	return flushSentry, nil
}

func StartSpan(ctx context.Context, name string, opts ...SpanOptionFunc) Span {
	return tracer.StartSpan(ctx, name, opts...)
}

func StartSpanWithContext(ctx context.Context, name string, opts ...SpanOptionFunc) (Span, context.Context) {
	sp := tracer.StartSpan(ctx, name, opts...)

	return sp, sp.Context()
}

// default flush time is 1 second
func flushSentry(flushTime string) {
	timeout, err := time.ParseDuration(flushTime + "s")

	if err != nil {
		timeout = 1 * time.Second
	}

	sentry.Flush(timeout)
}
