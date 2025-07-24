package tracer

import (
	"context"
	"github.com/getsentry/sentry-go"
	"math/rand"
)

type SentryTracer struct {
	isActive          bool
	successSampleRate float64
	rand              *rand.Rand
}

func (s *SentryTracer) StartSpan(ctx context.Context, name string, opts ...SpanOptionFunc) Span {
	// manage option
	options := &SpanOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// validate if tracer is not active
	if !s.isActive {
		return &NoopSpan{ctx: ctx}
	}

	// set logic to implement sample rate
	if !options.ForceRecord && s.rand.Float64() > s.successSampleRate {
		return &NoopSpan{ctx: ctx}
	}

	// initiate span
	var sp *sentry.Span

	// check active span from context
	if span := sentry.SpanFromContext(ctx); span == nil || span.Op == "" {
		sp = sentry.StartTransaction(ctx, name)

		// set default span data
		sp.Source = sentry.SourceURL
		sp.Origin = sentry.SpanOriginFiber
	} else {
		sp = sentry.StartSpan(ctx, name)
	}

	// data
	if sp.Data == nil {
		sp.Data = make(map[string]interface{})
	}

	// request option
	if options.Request != nil {
		sp.Data[KeyRequest] = options.Request
	}

	// attributes option
	if options.Attributes != nil && len(options.Attributes) != 0 {
		for key, val := range options.Attributes {
			if key == "operation_name" {
				if OpName, ok := val.(string); ok {
					sp.Op = OpName
				}
			} else {
				sp.Data[key] = val
			}
		}
	}

	return &SentrySpan{span: sp}
}

func (s *SentryTracer) StartSpanWithContext(ctx context.Context, name string, opts ...SpanOptionFunc) (Span, context.Context) {
	sp := tracer.StartSpan(ctx, name, opts...)

	return sp, sp.Context()
}
