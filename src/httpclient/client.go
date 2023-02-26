package httpclient

import (
	"errors"
	"github.com/sghaida/go-stuff/src/cauth"
	"net/http"
)

// Client ...
type Client struct {
	config   *Config
	client   http.Client
	authType cauth.IAuth
}

// NewClient create new http Client
func NewClient(config *Config, client *http.Client, authType cauth.IAuth) (*Client, error) {
	// set up the transport layer
	// allow 100 concurrent connection in the connection pool
	if client.Transport == nil {
		t := http.Transport{}
		client.Transport = &t
	}
	t := client.Transport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	// override transport
	client.Transport = t

	if config == nil {
		return nil, errors.New("config is empty")
	}
	if authType == nil {
		return nil, errors.New("auth type is not defined")
	}
	return &Client{client: *client, config: config, authType: authType}, nil
}

func (c *Client) getAuthHeader() (cauth.AuthHeader, error) {
	switch authData := c.authType.(type) {
	case cauth.ISomeAuth:
		return authData.GetAuthData()
	case cauth.INoAuth:
		return cauth.NewAuthHeader("", ""), nil
	default:
		return cauth.NewAuthHeader("", ""), errors.New("unsupported auth type")
	}
}
