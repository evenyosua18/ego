package tracer

import (
	"context"
	"github.com/evenyosua18/ego/tracer"
)

type StubTracer struct{}

func (s *StubTracer) StartSpan(ctx context.Context, name string, opts ...tracer.SpanOptionFunc) tracer.Span {
	return &StubSpan{
		Ctx:     ctx,
		TraceID: "TEST",
	}
}

func (s *StubTracer) StartSpanWithContext(ctx context.Context, name string, opts ...tracer.SpanOptionFunc) (tracer.Span, context.Context) {
	span := s.StartSpan(ctx, name, opts...)
	return span, span.Context()
}
