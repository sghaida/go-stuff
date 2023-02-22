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
	method  HttpMethod
	headers map[string]string
	query   map[string]string
	reqBody []byte
	client  *Client
}

// Call : do request http call with background context
// and if timeout is defined in the config, set the context timeout and call
func (c *Caller) Call() (*http.Response, error) {
	ctx := context.Background()
	return c.CallWithContext(ctx)
}

// CallWithContext do request http call with context
func (c *Caller) CallWithContext(ctx context.Context) (*http.Response, error) {

	// read request body
	body := io.NopCloser(bytes.NewReader(c.reqBody))

	// create the http request
	url := fmt.Sprintf("%s/%s", c.host, c.route)
	req, err := http.NewRequest(string(c.method), url, body)

	req = req.WithContext(ctx)

	// add the default headers  (from the config) if available
	for key, value := range c.client.config.defaultHeaders {
		req.Header.Add(key, value)
	}
	// add extra headers passed by the request
	for key, value := range c.headers {
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
	for k, v := range c.query {
		q.Add(k, v)
	}

	resp, err := c.client.client.Do(req)

	return resp, err
}

// RetryableCall do http call with retry logic.
func (c *Caller) RetryableCall() (*http.Response, error) {

	retryable := retry.NewRetry(c.client.config.numOfRetries, retry.DefaultInitialDelay, retry.DefaultMaxDelay)
	toExecute := func() (interface{}, error) {
		return c.Call()
	}
	resp, err := retryable.Run(toExecute)

	response, _ := resp.(*http.Response)
	return response, err

}
