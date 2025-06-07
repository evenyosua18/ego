package tracer

type (
	SpanOptionFunc func(*SpanOptions)

	SpanOptions struct {
		Attributes map[string]any
		Request    any
	}
)

func WithRequest(req any) SpanOptionFunc {
	return func(options *SpanOptions) {
		options.Request = req
	}
}

func WithAttributes(attributes map[string]any) SpanOptionFunc {
	return func(options *SpanOptions) {
		if options.Attributes == nil {
			options.Attributes = make(map[string]any)
		}

		for key, val := range attributes {
			options.Attributes[key] = val
		}
	}
}

func WithAttribute(key string, value any) SpanOptionFunc {
	return func(options *SpanOptions) {
		if options.Attributes == nil {
			options.Attributes = make(map[string]any)
		}

		options.Attributes[key] = value
	}
}

func WithFilePath(val string) SpanOptionFunc {
	return func(options *SpanOptions) {
		if options.Attributes == nil {
			options.Attributes = make(map[string]any)
		}

		options.Attributes["file_path"] = val
	}
}
