package main

import (
	"context"
	"errors"
	"fmt"
	"hit"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type env struct {
	stdout io.Writer // stdout abstracts standard output
	stderr io.Writer // stderr abstracts standard error
	args   []string  // args are command-line arguments
	dry    bool      // dry enables dry mode
}

func run(e *env) error {
	c := &config{
		n:      100,
		c:      1,
		method: GET,
	}
	if err := parseArgs(c, e.args, e.stderr); err != nil {
		fmt.Fprintf(e.stderr, "%s\n", err)
		return err
	}
	fmt.Fprintf(e.stdout, "%+v", *c)
	if e.dry {
		return nil
	}

	return runHit(e, c)
}

func runHit(e *env, c *config) error {

	handleErr := func(err error) error {
		if err != nil {
			fmt.Fprintf(e.stderr, "\n error occured: %v\n", err)
		}
		return err
	}

	const (
		timeout           = time.Hour
		timeoutPerRequest = 30 * time.Second
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	request, err := http.NewRequest(http.MethodGet, c.url, http.NoBody)
	if err != nil {
		return handleErr(fmt.Errorf("new request: %w", err))
	}

	client := &hit.Client{
		C:       c.c,
		RPS:     c.rps,
		Timeout: 30 * time.Second,
	}
	sum := client.Do(ctx, request, c.n)
	sum.Fprint(e.stdout)

	if err = ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		return handleErr(fmt.Errorf("timed out in %s", timeout))
	}
	return handleErr(err)
}
