package hit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Option allows changing Client's behavior.
type Option func(c *Client)

// Concurrency changes the Client's concurrency level.
func Concurrency(n int) Option {
	return func(c *Client) { c.C = n }
}

// RequestsPerSecond changes Client's RPS (requests persecond)
func RequestsPerSecond(d int) Option {
	return func(c *Client) { c.RPS = d }
}

// Timeout changes the Client's timeout per request
func Timeout(d time.Duration) Option {
	return func(c *Client) { c.Timeout = d }
}

// SendN sends n HTTP requests to the url and returns an
// aggregated [Result].
func SendN(ctx context.Context, url string, n int, opts ...Option) (Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return Result{}, fmt.Errorf("new http request: %w",
			err)
	}
	var c Client
	for _, opt := range opts {
		opt(&c)
	}
	return c.Do(ctx, r, n), nil
}
