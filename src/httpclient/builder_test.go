package httpclient

import (
	"github.com/sghaida/go-stuff/src/cauth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewCallerBuilder(t *testing.T) {
	tt := []struct {
		name          string
		host          string
		route         string
		method        HttpMethod
		query         map[string]string
		headers       map[string]string
		body          []byte
		configBuilder *ConfigBuilder
		buildClient   bool
		expectsError  bool
	}{
		{
			name:          "get without any param",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        GET,
			query:         map[string]string{},
			headers:       map[string]string{},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  false,
		},
		{
			name:          "Get with header and query params",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        GET,
			query:         map[string]string{"userId": "123"},
			headers:       map[string]string{"x-client-id": "123"},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  false,
		},
		{
			name:          "post with post body",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        POST,
			query:         map[string]string{},
			headers:       map[string]string{"x-client-id": "123"},
			body:          []byte(`{"user": "sghaida"}`),
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  false,
		},
		{
			name:          "put with body",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        PUT,
			query:         map[string]string{},
			headers:       map[string]string{"x-client-id": "123"},
			body:          []byte(`{"user": "sghaida"}`),
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  false,
		},
		{
			name:          "without host",
			host:          "",
			route:         "api/v1/users",
			method:        GET,
			query:         map[string]string{},
			headers:       map[string]string{},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  true,
		},
		{
			name:          "without route",
			host:          "https://example.com",
			route:         "",
			method:        GET,
			query:         map[string]string{},
			headers:       map[string]string{},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  true,
		},
		{
			name:          "without client",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        GET,
			query:         map[string]string{},
			headers:       map[string]string{},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   false,
			expectsError:  true,
		},
		{
			name:          "without method",
			host:          "https://example.com",
			route:         "api/v1/users",
			method:        "",
			query:         map[string]string{},
			headers:       map[string]string{},
			body:          nil,
			configBuilder: NewConfig().WithTimeout(1 * time.Second).WithRetry(3),
			buildClient:   true,
			expectsError:  true,
		},
	}

	for _, tc := range tt {

		config, _ := tc.configBuilder.Build()
		client, _ := NewClient(config, http.DefaultClient, &cauth.NoAuth{})

		t.Run(tc.name, func(t *testing.T) {
			var req *CallerBuilder
			if tc.buildClient {
				req = NewCallerBuilder(client, tc.host, tc.route, tc.method).
					WithHeaders(tc.headers).
					WithRequestBody(tc.body).
					WithQueryParam(tc.query)
			} else {
				req = NewCallerBuilder(nil, tc.host, tc.route, tc.method).
					WithHeaders(tc.headers).
					WithRequestBody(tc.body).
					WithQueryParam(tc.query)
			}

			caller, err := req.Build()

			if !tc.expectsError {
				assert.NoError(t, err)
				assert.Equal(t, caller.route, tc.route)
				assert.Equal(t, caller.host, tc.host)
				assert.Equal(t, caller.query, tc.query)
				assert.Equal(t, caller.headers, tc.headers)
				assert.Equal(t, caller.reqBody, tc.body)
				return
			}
			assert.Error(t, err)
		})
	}
}
