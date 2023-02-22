package httpclient

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	tt := []struct {
		name         string
		timeout      time.Duration
		retries      int
		expectsError bool
	}{
		{
			name:         "create successful config",
			timeout:      0 * time.Second,
			retries:      0,
			expectsError: false,
		},
		{
			name:         "negative timeout",
			timeout:      -1 * time.Second,
			retries:      0,
			expectsError: true,
		},
		{
			name:         "negative retries",
			timeout:      0 * time.Second,
			retries:      -1,
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cb, err := NewConfig().
				WithTimeout(tc.timeout).
				WithHeaders(map[string]string{"X-CLIENT_ID": "bla-bla-bla"}).
				WithRetry(tc.retries).
				WithJsonSchema([]byte("{}")).Build()
			if tc.expectsError {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, cb.timeout, tc.timeout)
			assert.Equal(t, cb.defaultHeaders["X-CLIENT_ID"], "bla-bla-bla")
			assert.Equal(t, cb.numOfRetries, tc.retries)
			assert.Equal(t, cb.jsonSchema, json.RawMessage("{}"))
		})

	}

}
