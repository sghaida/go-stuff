package cauth

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// BasicAuth ...
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth ...
func NewBasicAuth(username string, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

// GetAuthType ...
func (b *BasicAuth) GetAuthType() AuthType {
	return AuthBasic
}

// GetAuthData ...
func (b *BasicAuth) GetAuthData() (AuthHeader, error) {
	if b.username == "" || b.password == "" {
		return AuthHeader{}, errors.New("username and password can't be empty")
	}
	auth := fmt.Sprintf("%s:%s", b.username, b.password)
	encodedAuthStr := base64.StdEncoding.EncodeToString([]byte(auth))

	return AuthHeader{key: "Authorization", value: fmt.Sprintf("Basic %s", encodedAuthStr)}, nil
}
