package httpclient

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Config holds the HttpCaller config
type Config struct {
	host           string
	route          string
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
	host           string
	route          string
	timeout        time.Duration
	numOfRetries   int
	jsonSchema     json.RawMessage
	defaultHeaders map[string]string
}

func (c *ConfigBuilder) Build() (*Config, error) {
	if c.host == "" {
		return nil, errors.New("host can't be empty")
	}
	if c.route == "" {
		return nil, errors.New("route can't be empty")
	}
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

func (c *ConfigBuilder) WithHost(host string) *ConfigBuilder {
	// remove trailing slashes and spaces/
	host = strings.TrimSpace(host)
	host = strings.TrimSuffix(host, "/")

	c.host = host
	return c
}

func (c *ConfigBuilder) WithRoute(route string) *ConfigBuilder {
	// remove trailing slashes and spaces/
	route = strings.TrimSpace(route)
	route = strings.TrimSuffix(route, "/")

	c.route = route
	return c
}

func (c *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	c.timeout = timeout
	return c
}

func (c *ConfigBuilder) WithRetry(numOfRetries int) *ConfigBuilder {
	c.numOfRetries = numOfRetries
	return c
}

func (c *ConfigBuilder) WithJsonSchema(schema json.RawMessage) *ConfigBuilder {
	c.jsonSchema = schema
	return c
}

func (c *ConfigBuilder) WithHeaders(headers map[string]string) *ConfigBuilder {
	c.defaultHeaders = headers
	return c
}
