package tracer

import "context"

type StubSpan struct {
	Ctx     context.Context
	TraceID string
}

func (s *StubSpan) Context() context.Context {
	// return defined context
	return s.Ctx
}

func (s *StubSpan) End() {
	// do nothing
}

func (s *StubSpan) LogError(err error) error {
	// return the given error
	return err
}

func (s *StubSpan) Log(name string, obj any) {
	// do nothing
}

func (s *StubSpan) LogResponse(obj any) {
	// do nothing
}

func (s *StubSpan) GetTraceID() string {
	// return defined trace id
	return s.TraceID
}
