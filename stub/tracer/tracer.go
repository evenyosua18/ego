package tracer

import (
	"context"
	"github.com/evenyosua18/ego/tracer"
)

type StubTracer struct {
	SpanToReturn tracer.Span
}

func (s *StubTracer) StartSpan(ctx context.Context, name string, opts ...tracer.SpanOptionFunc) tracer.Span {
	if s.SpanToReturn != nil {
		return s.SpanToReturn
	}

	return &StubSpan{}
}

func (s *StubTracer) StartSpanWithContext(ctx context.Context, name string, opts ...tracer.SpanOptionFunc) (tracer.Span, context.Context) {
	span := s.StartSpan(ctx, name, opts...)
	return span, span.Context()
}
