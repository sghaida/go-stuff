package httpclient

import (
	"github.com/sghaida/go-stuff/src/cauth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tt := []struct {
		name         string
		numOfRetries int
		timeout      time.Duration
		auth         cauth.IAuth
		expectsError bool
	}{
		{
			name:         "create client successfully",
			numOfRetries: 0,
			timeout:      0,
			auth:         cauth.NoAuth,
			expectsError: false,
		},
		{
			name:         "basic auth client",
			numOfRetries: 0,
			timeout:      0,
			auth:         cauth.NewBasicAuth("user", "password"),
			expectsError: false,
		},
		{
			name:         "jwt auth client",
			numOfRetries: 0,
			timeout:      0,
			auth:         cauth.NewJWTAuth("some-token"),
			expectsError: false,
		},
		{
			name:         "api key client",
			numOfRetries: 0,
			timeout:      0,
			auth:         cauth.NewAPIKey("some-api-key"),
			expectsError: false,
		},
		{
			name:         "invalid auth",
			numOfRetries: 0,
			timeout:      0,
			auth:         nil,
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			config, _ := NewConfig().WithTimeout(tc.timeout).WithRetry(tc.numOfRetries).Build()
			client, err := NewClient(config, http.DefaultClient, tc.auth)
			if tc.expectsError {
				assert.Error(t, err)
				return
			}
			assert.NotNil(t, client)
			auth, err := client.getAuthHeader()
			if err != nil {
				assert.Failf(t, "expected auth header, got error: %v", err.Error())
			}
			key, value := auth.GetAuthKeyValue()
			authType := tc.auth.GetAuthType()
			if authType != cauth.None {
				assert.True(t, key != "")
				assert.True(t, value != "")
			}

		})
	}
}
