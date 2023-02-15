package httpclient

import "net/http"

// HttpCaller interface provides the definition of the caller method
type HttpCaller interface {
	// Call : executes http request
	Call(method HttpMethod, headers map[string]string, query map[string]string, reqBody []byte) *http.Response
	// CallWithContext : executes http request with context
	CallWithContext(
		method HttpMethod, headers map[string]string, query map[string]string, reqBody []byte,
	) (*http.Response, error)
	// RetryableCall executes Call function in a retryable manner
	RetryableCall(
		method HttpMethod, extraHeaders map[string]string, query map[string]string, reqBody []byte,
	) (*http.Response, error)
}
