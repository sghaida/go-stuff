package httpclient

import "errors"

// CallerBuilder ...
type CallerBuilder struct {
	host    string
	route   string
	method  HttpMethod
	headers map[string]string
	query   map[string]string
	reqBody []byte
	client  *Client
}

// NewCallerBuilder creates http CallerBuilder
func NewCallerBuilder(client *Client, host, route string, method HttpMethod) *CallerBuilder {
	return &CallerBuilder{
		host:    host,
		route:   route,
		method:  method,
		headers: make(map[string]string),
		query:   make(map[string]string),
		client:  client,
	}
}

// WithHeaders add request headers
func (b *CallerBuilder) WithHeaders(headers map[string]string) *CallerBuilder {
	if len(headers) != 0 {
		for k, v := range headers {
			b.headers[k] = v
		}
	}
	return b
}

// WithQueryParam add request query param
func (b *CallerBuilder) WithQueryParam(params map[string]string) *CallerBuilder {
	if len(params) != 0 {
		for k, v := range params {
			b.query[k] = v
		}
	}
	return b
}

// WithRequestBody add requestBody
func (b *CallerBuilder) WithRequestBody(reqBody []byte) *CallerBuilder {
	if len(reqBody) != 0 {
		b.reqBody = reqBody
	}
	return b
}

// Build : Build http Caller
func (b *CallerBuilder) Build() (*Caller, error) {
	if b.client == nil {
		return nil, errors.New("client can't be nil")
	}
	if b.route == "" {
		return nil, errors.New("http route can't be empty")
	}
	if b.host == "" {
		return nil, errors.New("http host can't be empty")
	}
	if b.method == "" {
		return nil, errors.New("http method can't be empty")
	}

	caller := &Caller{
		host:    b.host,
		route:   b.route,
		method:  b.method,
		headers: b.headers,
		query:   b.query,
		reqBody: b.reqBody,
		client:  b.client,
	}
	return caller, nil
}
