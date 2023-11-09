package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	jsonEncoder = func(v interface{}) (io.Reader, error) {
		var buf bytes.Buffer

		return &buf, json.NewEncoder(&buf).Encode(v)
	}
	jsonHeader = String("Content-Type", "application/json")
)

var ErrBody = errors.New("must set []byte, io.Reader, or string for the body")

type Request struct {
	base         *url.URL
	path         string
	args         []interface{}
	method       string
	query        url.Values
	headers      []Value
	body         interface{}
	encoder      func(v interface{}) (io.Reader, error)
	errorHandler func(statusCode int, b []byte) error
}

func NewRequest(method string) Request {
	return Request{
		method: method,
	}
}

func GET(handler string, args ...interface{}) Request {
	return NewRequest(http.MethodGet).WithPathf(handler, args...)
}

func POST(handler string, args ...interface{}) Request {
	return NewRequest(http.MethodPost).WithPathf(handler, args...)
}

func PUT(handler string, args ...interface{}) Request {
	return NewRequest(http.MethodPut).WithPathf(handler, args...)
}

func DELETE(handler string, args ...interface{}) Request {
	return NewRequest(http.MethodDelete).WithPathf(handler, args...)
}

func (r Request) WithPathf(path string, args ...interface{}) Request {
	r.path = path
	r.args = args

	return r
}

func (r Request) WithPath(path string) Request {
	r.path = path

	return r
}

func (r Request) Method() string {
	return r.method
}

func (r Request) WithBase(u *url.URL) Request {
	r.base = u

	return r
}

func (r Request) Handler() string {
	return r.path
}

func (r Request) URI() string {
	u := r.path
	if len(r.args) > 0 {
		u = fmt.Sprintf(r.path, r.args...)
	}

	if len(r.query) > 0 {
		return u + "?" + r.query.Encode()
	}

	return u
}

func (r Request) WithQuery(value ...Value) Request {
	if r.query == nil {
		r.query = make(url.Values, len(value))
	}

	for _, v := range value {
		v(r.query)
	}

	return r
}

func (r Request) WithHeader(value ...Value) Request {
	r.headers = append(r.headers, value...)

	return r
}

func (r Request) WithJSONBody(b interface{}) Request {
	return r.WithBody(b, jsonEncoder, jsonHeader)
}

func (r Request) WithBody(
	b interface{},
	encoder func(v interface{}) (io.Reader, error),
	headers ...Value,
) Request {
	r.body = b
	r.encoder = encoder
	r.headers = append(r.headers, headers...)

	return r
}

func (r Request) RawBody() interface{} {
	return r.body
}

func (r Request) Body() (io.Reader, error) {
	if r.body == nil {
		return nil, nil
	}

	switch data := r.body.(type) {
	case string:
		return bytes.NewBufferString(data), nil
	case []byte:
		return bytes.NewBuffer(data), nil
	case io.Reader:
		return data, nil
	default:
		if r.encoder != nil {
			return r.encoder(r.body)
		}

		return nil, ErrBody
	}
}

func (r Request) ParseURL() (*url.URL, error) {
	u, err := r.base.Parse(r.URI())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid url", err)
	}

	return u, nil
}

func (r Request) HTTP(ctx context.Context) (*http.Request, error) {
	u, err := r.ParseURL()
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	b, err := r.Body()
	if err != nil {
		return nil, fmt.Errorf("create body error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, r.method, u.String(), b)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	for _, v := range r.headers {
		v(req.Header)
	}

	return req, nil
}

func (r Request) WithErrorHandler(h func(statusCode int, b []byte) error) {
	r.errorHandler = h
}
