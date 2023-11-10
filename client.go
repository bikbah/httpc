package httpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 15 * time.Second

type Client struct {
	client  *http.Client
	url     *url.URL
	name    string
	decode  func(data []byte, v interface{}) error
	headers []Value
	logFunc func(
		ctx context.Context,
		name string,
		method string,
		url string,
		dur time.Duration,
		statusCode int,
		v any,
		b []byte,
		err error,
	)
}

func New(u *url.URL, opts ...Option) *Client {
	c := &Client{
		url: u,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		name:   resolveName(u),
		decode: json.Unmarshal,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.client.Transport == nil {
		c.client.Transport = http.DefaultTransport
	}

	return c
}

func Must(base string, opts ...Option) *Client {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	return New(u, opts...)
}

type Option func(*Client)

func WithTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.client.Timeout = t
	}
}

func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) {
		c.client.Transport = t
	}
}

func WithLogFunc(logFunc func(
	ctx context.Context,
	name string,
	method string,
	url string,
	dur time.Duration,
	statusCode int,
	v any,
	b []byte,
	err error,
),
) Option {
	return func(c *Client) {
		c.logFunc = logFunc
	}
}

func WithDecode(d func(data []byte, v interface{}) error) Option {
	return func(c *Client) {
		c.decode = d
	}
}

func WithHeaders(h ...Value) Option {
	return func(c *Client) {
		c.headers = h
	}
}

func WithName(n string) Option {
	return func(c *Client) {
		c.name = n
	}
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Do(ctx context.Context, r Request, v interface{}, customParseFunc func(resp *http.Response, b []byte) error) (err error) {
	r = r.WithBase(c.url)

	start := time.Now().UTC()
	var (
		url              string
		rawResponseBytes []byte
		statusCode       int
	)
	defer func() {
		if c.logFunc != nil {
			c.logFunc(
				ctx,
				c.Name(),
				r.Method(),
				url,
				time.Since(start),
				statusCode,
				v,
				rawResponseBytes,
				err,
			)
		}
	}()

	req, err := r.WithHeader(c.headers...).HTTP(ctx)
	if err != nil {
		return fmt.Errorf("%s: create http request error: %w", c.Name(), err)
	}

	url = req.URL.String()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: do request error: %w", c.Name(), err)
	}
	defer resp.Body.Close()

	rawResponseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: read body error: %w", c.Name(), err)
	}

	statusCode = resp.StatusCode
	if customParseFunc != nil {
		return customParseFunc(resp, rawResponseBytes)
	}

	if resp.StatusCode >= 400 {
		if r.errorHandler != nil {
			return r.errorHandler(resp.StatusCode, rawResponseBytes)
		}

		return fmt.Errorf("%s: http status code: %v", c.Name(), resp.StatusCode)
	}

	if len(rawResponseBytes) > 1 && c.decode != nil && v != nil {
		return c.decode(rawResponseBytes, v)
	}

	return nil
}

func resolveName(u *url.URL) string {
	if idx := strings.Index(u.Hostname(), "."); idx > 0 {
		return strings.ToLower(u.Hostname()[:idx])
	}

	return strings.ToLower(u.Hostname())
}
