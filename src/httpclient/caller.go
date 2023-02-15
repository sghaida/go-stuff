package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sghaida/go-stuff/src/cauth"
	"github.com/sghaida/go-stuff/src/retry"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
)

// Caller ...
type Caller struct {
	config     *Config
	authHeader map[string]string
	client     http.Client
}

// NewHTTPCaller create new http caller
func NewHTTPCaller(config *Config) *Caller {
	header := make(map[string]string)
	// set up the transport layer
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	client := http.Client{
		Timeout:   config.timeout * time.Second,
		Transport: t,
	}

	return &Caller{client: client, config: config, authHeader: header}

}

func (c *Caller) WithAuth(authType cauth.IAuth) (*Caller, error) {
	// check if auth header is already being set
	if len(c.authHeader) > 0 {
		return c, errors.New("authentication is already being set")
	}
	switch authData := authType.(type) {
	case cauth.IBasicAuth:
		return c.addBasicAuth(authData)
	case cauth.IJwtAuth:
		return c.addJwtToken(authData)
	case cauth.IApiKey:
		return c.addApiKey(authData)
	case cauth.INoAuth:
		return c, nil
	default:
		return nil, errors.New("unsupported auth type")
	}
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
	// remove trailing / character
	host := c.config.host
	path := c.config.route

	if host == "" {
		return nil, errors.New("http host can't be empty")
	}

	if method == "" {
		return nil, errors.New("http method can't be empty")
	}
	// remove trailing slashes and spaces/
	host = strings.TrimSpace(host)
	host = strings.TrimSuffix(host, "/")
	path = strings.TrimSpace(path)
	path = strings.TrimSuffix(path, "/")

	// read request body
	body := io.NopCloser(bytes.NewReader(reqBody))

	// create the http request
	url := fmt.Sprintf("%s/%s", host, path)
	req, err := http.NewRequest(string(method), url, body)

	req = req.WithContext(ctx)

	// add the default headers  (from the config) if available
	for key, value := range c.config.defaultHeaders {
		req.Header.Add(key, value)
	}
	// add extra headers passed by the request
	for key, value := range extraHeaders {
		req.Header.Add(key, value)
	}
	// add auth header
	for key, value := range c.authHeader {
		req.Header.Add(key, value)
	}
	// add query values
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	resp, err := c.client.Do(req)

	return resp, err
}

// RetryableCall do http call with retry logic.
func (c *Caller) RetryableCall(
	method HttpMethod, extraHeaders map[string]string, query map[string]string, reqBody []byte,
) (*http.Response, error) {

	retryable := retry.NewRetry(c.config.numOfRetries, retry.DefaultInitialDelay, retry.DefaultMaxDelay)
	toExecute := func() (interface{}, error) {
		return c.Call(method, extraHeaders, query, reqBody)
	}
	resp, err := retryable.Run(toExecute)

	response, _ := resp.(*http.Response)
	return response, err

}
