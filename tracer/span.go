package tracer

import (
	"context"
	"github.com/getsentry/sentry-go"
)

const (
	KeyResponse = "response"
	KeyError    = "error"
	KeyRequest  = "request"
)

type SentrySpan struct {
	span *sentry.Span
}

func (sp *SentrySpan) Context() context.Context {
	return sp.span.Context()
}

func (sp *SentrySpan) End() {
	if sp.span.Status == sentry.SpanStatusUnknown {
		sp.span.Status = sentry.SpanStatusOK
	}

	sp.span.Finish()
}

func (sp *SentrySpan) LogError(err error) error {
	sp.span.Status = sentry.SpanStatusInternalError
	sp.span.Data[KeyError] = err

	return err
}

func (sp *SentrySpan) Log(name string, obj any) {
	sp.span.Data[name] = obj
}

func (sp *SentrySpan) LogResponse(obj any) {
	sp.span.Data[KeyResponse] = obj
}

func (sp *SentrySpan) GetTraceID() string {
	return sp.span.TraceID.String()
}

// NoopSpan no op span for default
type NoopSpan struct {
	ctx context.Context
}

func (sp *NoopSpan) Context() context.Context {
	return sp.ctx
}

func (sp *NoopSpan) End() {}

func (sp *NoopSpan) LogError(err error) error {
	return err
}

func (sp *NoopSpan) Log(name string, obj any) {}

func (sp *NoopSpan) LogResponse(obj any) {}

func (sp *NoopSpan) GetTraceID() string {
	return ""
}
