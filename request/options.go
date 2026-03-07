package request

import "time"

// RequestOption is a functional option for configuring a Request.
type RequestOption func(*Request)

// WithHeaders sets custom headers for the request.
// It will overwrite any existing header with the same key.
func WithHeaders(headers map[string]string) RequestOption {
	return func(r *Request) {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
}

// WithQueryParams sets query parameters for the request URL.
func WithQueryParams(params map[string]string) RequestOption {
	return func(r *Request) {
		if r.QueryParams == nil {
			r.QueryParams = make(map[string]string)
		}
		for k, v := range params {
			r.QueryParams[k] = v
		}
	}
}

// WithBody sets the request body. It can be a byte slice or a struct to be marshaled to JSON.
func WithBody(body any) RequestOption {
	return func(r *Request) {
		r.Body = body
	}
}

// WithRetry enables retrying the request upon specific failures (network errors, 429, or user defined status codes).
func WithRetry(maxRetries int, delay time.Duration, retryCodes ...int) RequestOption {
	return func(r *Request) {
		r.maxRetries = maxRetries
		r.retryDelay = delay
		r.retryCodes = retryCodes
	}
}
