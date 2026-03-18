package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/evenyosua18/ego/code"
	"github.com/sony/gobreaker/v2"
)

// Client is a wrapper around http.Client
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new HTTP wrapper client.
// If httpClient is nil, http.DefaultClient is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
	}
}

// Request holds all necessary data to perform an HTTP request.
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        any

	// Retry configuration
	maxRetries int
	retryDelay time.Duration
	retryCodes []int

	// Circuit breaker bypass
	disableBreaker bool
}

// applyOptions processes functional options for the request.
func (r *Request) applyOptions(opts ...RequestOption) {
	for _, opt := range opts {
		opt(r)
	}
}

// Get executes an HTTP GET request.
func (c *Client) Get(ctx context.Context, url string, opts ...RequestOption) (*http.Response, []byte, error) {
	req := Request{Method: http.MethodGet, URL: url}
	req.applyOptions(opts...)
	return c.Do(ctx, req)
}

// Post executes an HTTP POST request.
func (c *Client) Post(ctx context.Context, url string, opts ...RequestOption) (*http.Response, []byte, error) {
	req := Request{Method: http.MethodPost, URL: url}
	req.applyOptions(opts...)
	return c.Do(ctx, req)
}

// Put executes an HTTP PUT request.
func (c *Client) Put(ctx context.Context, url string, opts ...RequestOption) (*http.Response, []byte, error) {
	req := Request{Method: http.MethodPut, URL: url}
	req.applyOptions(opts...)
	return c.Do(ctx, req)
}

// Delete executes an HTTP DELETE request.
func (c *Client) Delete(ctx context.Context, url string, opts ...RequestOption) (*http.Response, []byte, error) {
	req := Request{Method: http.MethodDelete, URL: url}
	req.applyOptions(opts...)
	return c.Do(ctx, req)
}

// Do executes the HTTP request with optional retries.
// Retries occur for network errors, timeouts, 5xx server errors, or 429 Too Many Requests responses.
func (c *Client) Do(ctx context.Context, req Request) (*http.Response, []byte, error) {
	var lastErr error
	var resp *http.Response
	var body []byte

	type attemptResult struct {
		resp *http.Response
		body []byte
		err  error
	}

	// Ensure at least 1 attempt (0 retries = 1 attempt)
	attempts := req.maxRetries + 1

	for i := 0; i < attempts; i++ {
		var res any
		var cbErr error

		if req.disableBreaker {
			r, b, e := c.doOnce(ctx, req)
			res = attemptResult{resp: r, body: b, err: e}
		} else {
			// Attempt the request via circuit breaker
			res, cbErr = executeWithBreaker(req.URL, func() (any, error) {
				r, b, e := c.doOnce(ctx, req)
				
				if e != nil {
					return attemptResult{resp: r, body: b, err: e}, e
				}
				
				if r != nil && (r.StatusCode >= 500 || r.StatusCode == http.StatusTooManyRequests) {
					return attemptResult{resp: r, body: b, err: nil}, errors.New("server error")
				}
				
				return attemptResult{resp: r, body: b, err: nil}, nil
			})
		}

		if res != nil {
			payload := res.(attemptResult)
			resp = payload.resp
			body = payload.body
			lastErr = payload.err
		}

		if cbErr != nil && (errors.Is(cbErr, gobreaker.ErrOpenState) || errors.Is(cbErr, gobreaker.ErrTooManyRequests)) {
			lastErr = code.Wrap(cbErr, code.InternalError)
			break // circuit breaker prevents further attempts
		}

		// Decide if we should retry
		shouldRetry := false

		// 1. Retry on network error or timeout
		if lastErr != nil {
			shouldRetry = true
		} else if resp != nil {
			// 2. Retry on specific HTTP status codes (by request, plus 429 Rate Limit)
			if resp.StatusCode == http.StatusTooManyRequests || slices.Contains(req.retryCodes, resp.StatusCode) {
				shouldRetry = true
			}
		}

		// If it succeeded or failed but shouldn't be retried, return immediately
		if !shouldRetry {
			return resp, body, lastErr
		}

		// Delay before the next retry (if it's not the last attempt)
		if i < attempts-1 {
			// Avoid busy waiting, use context to handle early cancellations during retry pauses
			select {
			case <-ctx.Done():
				// Context cancelled/timed out, stop retrying
				return resp, body, ctx.Err()
			case <-time.After(req.retryDelay):
				// Proceed to the next loop iteration (retry)
			}
		}
	}

	// Always return the result of the final attempt
	return resp, body, lastErr
}

// doOnce handles the concrete execution of a single HTTP request attempt.
func (c *Client) doOnce(ctx context.Context, req Request) (*http.Response, []byte, error) {
	reqURL, err := url.Parse(req.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if len(req.QueryParams) > 0 {
		q := reqURL.Query()
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		reqURL.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	if req.Body != nil {
		if b, ok := req.Body.([]byte); ok {
			bodyReader = bytes.NewReader(b)
		} else {
			jsonBody, err := json.Marshal(req.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal request body to JSON: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBody)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request: %w", err)
	}

	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	outReq, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("http client failed to execute request: %w", err)
	}
	defer outReq.Body.Close()

	respBody, err := io.ReadAll(outReq.Body)
	if err != nil {
		return outReq, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return outReq, respBody, nil
}
