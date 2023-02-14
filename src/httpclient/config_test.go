package httpclient

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cb, _ := NewConfig().
		WithHost("https://somehost.com/").
		WithRoute("/resource/").
		WithTimeout(1 * time.Second).
		WithHeaders(map[string]string{"X-CLIENT_ID": "bla-bla-bla"}).
		WithRetry(3).
		WithJsonSchema([]byte("{}")).Build()

	assert.Equal(t, cb.host, "https://somehost.com")
	assert.Equal(t, cb.route, "/resource")
	assert.Equal(t, cb.timeout, 1*time.Second)
	assert.Equal(t, cb.defaultHeaders["X-CLIENT_ID"], "bla-bla-bla")
	assert.Equal(t, cb.numOfRetries, 3)
	assert.Equal(t, cb.jsonSchema, json.RawMessage("{}"))
}
