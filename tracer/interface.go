package tracer

import "context"

type Tracer interface {
	StartSpan(ctx context.Context, name string, opts ...SpanOptionFunc) Span
}

type Span interface {
	Context() context.Context
	End()
	LogError(err error) error
	Log(name string, obj any)
	LogResponse(obj any)
	GetTraceID() string
}
