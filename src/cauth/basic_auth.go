package cauth

import (
	"encoding/base64"
	"errors"
	"fmt"
)

type BasicAuth struct {
	username string
	password string
}

func NewBasicAuth(username string, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

func (b *BasicAuth) GetAuthType() AuthType {
	return AuthBasic
}

func (b *BasicAuth) GetAuthKey() (string, error) {
	if b.username == "" || b.password == "" {
		return "", errors.New("username and password can't be empty")
	}
	auth := fmt.Sprintf("%s:%s", b.username, b.password)
	encodedAuthStr := base64.StdEncoding.EncodeToString([]byte(auth))
	return encodedAuthStr, nil
}
