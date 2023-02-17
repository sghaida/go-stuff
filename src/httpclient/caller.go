package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sghaida/go-stuff/src/retry"
	"io"
	"net/http"
)

// TODO 3: update the retryable logic to include error codes to retry
const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
)

type Caller struct {
	host    string
	route   string
	headers map[string]string
	query   map[string]string
	reqBody []byte
	client  *Client
}

// NewCaller creates HTTPCaller
func NewCaller(client *Client, host, route string) *Caller {
	return &Caller{
		host:    host,
		route:   route,
		headers: make(map[string]string),
		query:   make(map[string]string),
		client:  client,
	}
}

// WithHeaders add request headers
func (c *Caller) WithHeaders(headers map[string]string) *Caller {
	if len(headers) != 0 {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
	return c
}

// WithQueryParam add request query param
func (c *Caller) WithQueryParam(params map[string]string) *Caller {
	if len(params) != 0 {
		for k, v := range params {
			c.query[k] = v
		}
	}
	return c
}

// WithRequestBody add requestBody
func (c *Caller) WithRequestBody(reqBody []byte) *Caller {
	if len(reqBody) != 0 {
		c.reqBody = reqBody
	}
	return c
}

func (c *Caller) Build() (*Caller, error) {
	if c.client == nil {
		return nil, errors.New("client can't be nil")
	}
	if c.route == "" {
		return nil, errors.New("http route can't be empty")
	}
	if c.host == "" {
		return nil, errors.New("http host can't be empty")
	}
	return c, nil
}

// Call : do request http call with background context
// and if timeout is defined in the config, set the context timeout and call
func (c *Caller) Call(
	method HttpMethod, headers map[string]string, query map[string]string, reqBody []byte,
) (*http.Response, error) {
	ctx := context.Background()
	return c.CallWithContext(ctx, method, headers, query, reqBody)
}

// CallWithContext do request http call with context
func (c *Caller) CallWithContext(
	ctx context.Context, method HttpMethod, extraHeaders map[string]string, query map[string]string, reqBody []byte,
) (*http.Response, error) {

	if method == "" {
		return nil, errors.New("http method can't be empty")
	}

	// read request body
	body := io.NopCloser(bytes.NewReader(reqBody))

	// create the http request
	url := fmt.Sprintf("%s/%s", c.host, c.route)
	req, err := http.NewRequest(string(method), url, body)

	req = req.WithContext(ctx)

	// add the default headers  (from the config) if available
	for key, value := range c.client.config.defaultHeaders {
		req.Header.Add(key, value)
	}
	// add extra headers passed by the request
	for key, value := range extraHeaders {
		req.Header.Add(key, value)
	}
	auth, err := c.client.getAuthHeader()
	if err != nil {
		// TODO wrap the error
		return nil, errors.New("unable to extract auth header")
	}
	// add auth header
	key, value := auth.GetAuthKeyValue()
	// skip no-auth case
	if key != "" && value != "" {
		req.Header.Add(key, value)
	}
	// add query values
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	resp, err := c.client.client.Do(req)

	return resp, err
}

// RetryableCall do http call with retry logic.
func (c *Caller) RetryableCall(
	method HttpMethod, extraHeaders map[string]string, query map[string]string, reqBody []byte,
) (*http.Response, error) {

	retryable := retry.NewRetry(c.client.config.numOfRetries, retry.DefaultInitialDelay, retry.DefaultMaxDelay)
	toExecute := func() (interface{}, error) {
		return c.Call(method, extraHeaders, query, reqBody)
	}
	resp, err := retryable.Run(toExecute)

	response, _ := resp.(*http.Response)
	return response, err

}
