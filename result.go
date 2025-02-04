package hit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Result is an [http.Request]'s performance result.
// Its zero value is useful.
type Result struct {
	RPS      float64       // RPS is the requests per second
	Requests int           // Requests is the number of requests made
	Errors   int           // Errors is the number of errors occurred
	Bytes    int64         // Bytes is the number of bytes downloaded
	Duration time.Duration // Duration is a single or all requests' duration
	Fastest  time.Duration // Fastest request result duration among others
	Slowest  time.Duration // Slowest request result duration among others
	Status   int           // Status is a request's HTTP status code
	Error    error         // Error is not nil if the request is failed
}

// Merge merges this [Result] with another.
func (r Result) Merge(other Result) Result {
	r.Requests++
	r.Bytes += other.Bytes
	if other.Error != nil || other.Status != http.StatusOK {
		r.Errors++
	}
	if r.Fastest == 0 || other.Duration < r.Fastest {
		r.Fastest = other.Duration
	}
	if other.Duration > r.Slowest {
		r.Slowest = other.Duration
	}
	return r
}

// Finalize calculates the total duration and RPS.
func (r Result) Finalize(total time.Duration) Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

// Fprint writes [Result] to an [io.Writer].
func (r Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	p("\nSummary:\n")
	p("\tSuccess		: %.0f%%\n", r.successRatio())
	p("\tRPS		: %.1f\n", r.RPS)
	p("\tRequests	: %d\n", r.Requests)
	p("\tErrors		: %d\n", r.Errors)
	p("\tBytes		: %d\n", r.Bytes)
	p("\tDuration	: %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest		: %s\n", round(r.Fastest))
		p("\tSlowest		: %s\n", round(r.Slowest))
	}
}

func (r Result) successRatio() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}

// String returns result as string.
func (r Result) String() string {
	var s strings.Builder
	r.Fprint(&s)
	return s.String()
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Millisecond)
}
