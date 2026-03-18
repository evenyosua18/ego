package request

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/evenyosua18/ego/code"
	"github.com/sony/gobreaker/v2"
)

var (
	breakers sync.Map

	globalBreakerConfig = BreakerConfig{
		MaxRequests:         1,
		Interval:            0,
		Timeout:             30 * time.Second,
		ConsecutiveFailures: 5,
	}
)

type BreakerConfig struct {
	MaxRequests         uint32
	Interval            time.Duration
	Timeout             time.Duration
	ConsecutiveFailures uint32
}

func InitBreakerConfig(config BreakerConfig) {
	globalBreakerConfig = config
}

// getBreaker retrieves or creates a circuit breaker for a given hostname.
func getBreaker(hostname string) *gobreaker.CircuitBreaker[any] {
	if cb, ok := breakers.Load(hostname); ok {
		return cb.(*gobreaker.CircuitBreaker[any])
	}

	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        hostname,
		MaxRequests: globalBreakerConfig.MaxRequests,
		Interval:    globalBreakerConfig.Interval,
		Timeout:     globalBreakerConfig.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if there are more than predefined consecutive failures
			return counts.ConsecutiveFailures > globalBreakerConfig.ConsecutiveFailures
		},
	})

	actual, _ := breakers.LoadOrStore(hostname, cb)
	return actual.(*gobreaker.CircuitBreaker[any])
}

// CheckEndpoints verifies if the circuit breakers for the given URLs are closed or half-open.
// It returns an error if any of the associated endpoints currently have an open circuit breaker.
func CheckEndpoints(urls ...string) error {
	for _, rawURL := range urls {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			continue // skip invalid URLs
		}
		
		host := parsed.Hostname()
		if host == "" {
			continue
		}

		cb := getBreaker(host)
		if cb.State() == gobreaker.StateOpen {
			err := fmt.Errorf("circuit breaker is open for host: %s", host)
			return code.Wrap(err, code.InternalError)
		}
	}
	return nil
}

// executeWithBreaker executes the given function within the circuit breaker for the given URL.
func executeWithBreaker(reqURL string, fn func() (any, error)) (any, error) {
	parsed, err := url.Parse(reqURL)
	if err != nil || parsed.Hostname() == "" {
		// If we can't parse or extract a hostname, just run the function directly.
		return fn()
	}

	cb := getBreaker(parsed.Hostname())
	return cb.Execute(fn)
}
