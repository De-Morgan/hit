package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

type config struct {
	url    string
	n      int
	c      int
	rps    int
	method httpMethod
}

func parseArgs(c *config, args []string, stderr io.Writer) error {
	fs := flag.NewFlagSet("hit", flag.ExitOnError)
	fs.SetOutput(stderr)
	fs.Var(newPositiveIntValue(&c.n), "n", "Number of requests")
	fs.Var(newPositiveIntValue(&c.c), "c", "Concurrency level")
	fs.Var(newPositiveIntValue(&c.rps), "rps", "Requests per second")
	fs.Var(&c.method, "m", "HTTP method for the request")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: %s [options] url\n ", fs.Name())
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	c.url = fs.Arg(0)
	if err := validateArgs(c); err != nil {
		fmt.Fprintln(fs.Output(), err)
		fs.Usage()
		return err
	}
	return nil
}
func validateArgs(c *config) error {
	const urlArg = "url argument"
	u, err := url.Parse(c.url)
	if err != nil {
		return argError(c.url, urlArg, err)
	}
	if c.url == "" || u.Host == "" || u.Scheme == "" {
		return argError(c.url, urlArg, errors.New("require a valid url"))
	}
	if c.n < c.c {
		err = fmt.Errorf(`should be greated than -c: "%d"`, c.c)
		return argError(c.n, "flag -n", err)
	}
	return nil
}
func argError(value any, arg string, err error) error {
	return fmt.Errorf(`invalid value "%v" for %s: %w`, value, arg, err)
}

var _ flag.Value = new(positiveIntValue)

type positiveIntValue int

func newPositiveIntValue(p *int) *positiveIntValue {
	return (*positiveIntValue)(p)
}
func (n *positiveIntValue) String() string {
	return strconv.Itoa(int(*n))
}

func (n *positiveIntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	if v <= 0 {
		return errors.New("should be greater than zero")
	}
	*n = positiveIntValue(v)
	return nil
}

type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	DELETE httpMethod = "DELETE"
)

func (h *httpMethod) String() string {
	return string(*h)
}
func (h *httpMethod) Set(s string) error {
	switch s {
	case string(GET):
		*h = GET
	case string(POST):
		*h = POST
	case string(PUT):
		*h = PUT
	case string(DELETE):
		*h = DELETE
	default:
		return fmt.Errorf("must be one of %s, %s, %s or %s", GET, POST, PUT, DELETE)
	}
	return nil
}
