package httpclient

import (
	"encoding/json"
	"errors"
	"time"
)

// Config holds the HttpCaller config
type Config struct {
	timeout        time.Duration
	numOfRetries   int
	jsonSchema     json.RawMessage
	defaultHeaders map[string]string
}

func newConfig(c *ConfigBuilder) *Config {
	config := (*Config)(c)
	return config
}

// ConfigBuilder Build HttpCaller config
type ConfigBuilder struct {
	timeout        time.Duration
	numOfRetries   int
	jsonSchema     json.RawMessage
	defaultHeaders map[string]string
}

// Build builds HttpCaller Config
func (c *ConfigBuilder) Build() (*Config, error) {
	if c.numOfRetries < 0 {
		return nil, errors.New("retries can't be negative")
	}
	if c.timeout < 0 {
		return nil, errors.New("timeout can't be negative")
	}

	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	if c.jsonSchema == nil {
		c.jsonSchema = json.RawMessage("{}")
	}

	conf := newConfig(c)
	return conf, nil
}

// NewConfig create new HttpCaller config
func NewConfig() *ConfigBuilder {
	return &ConfigBuilder{}
}

// WithTimeout add timeout
func (c *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	c.timeout = timeout
	return c
}

// WithRetry add retry count
func (c *ConfigBuilder) WithRetry(numOfRetries int) *ConfigBuilder {
	c.numOfRetries = numOfRetries
	return c
}

// WithJsonSchema add json schema
func (c *ConfigBuilder) WithJsonSchema(schema json.RawMessage) *ConfigBuilder {
	c.jsonSchema = schema
	return c
}

// WithHeaders add http headers
func (c *ConfigBuilder) WithHeaders(headers map[string]string) *ConfigBuilder {
	c.defaultHeaders = headers
	return c
}
