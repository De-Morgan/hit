package hit

import (
	"context"
	"io"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	C       int // C is the concurrency level
	RPS     int // RPS throttles the requests per second
	Timeout time.Duration
}

// Do sends n HTTP requests and returns an aggregated result.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) Result {
	now := time.Now()
	var sum Result
	for result := range c.do(ctx, r, n) {
		sum = sum.Merge(result)
	}
	return sum.Finalize(time.Since(now))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) <-chan Result {
	pipe := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})
	if c.RPS > 0 {
		t := time.Second / time.Duration(c.RPS*c.concurrency())
		pipe = throttle(pipe, t)
	}
	client := c.client()
	defer client.CloseIdleConnections()

	return split(pipe, c.concurrency(), func(r *http.Request) Result {
		result, _ := Send(client, r)
		return result
	})

}
func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.concurrency(),
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

}
func Send(c *http.Client, r *http.Request) (Result, error) {
	t := time.Now()
	var (
		code  int
		bytes int64
	)
	response, err := c.Do(r)
	if err == nil { // no err
		code = response.StatusCode
		bytes, err = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}
	result := Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
	return result, err
}
