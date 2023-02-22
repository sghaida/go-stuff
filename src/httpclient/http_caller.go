package httpclient

import (
	"context"
	"net/http"
)

// HttpCaller interface provides the definition of the Client method
type HttpCaller interface {
	// Call : executes http request
	Call() *http.Response
	// CallWithContext : executes http request with context
	CallWithContext(ctx context.Context) (*http.Response, error)
	// RetryableCall executes Call function in a retryable manner
	RetryableCall() (*http.Response, error)
}
