package cauth

import (
	"errors"
	"fmt"
)

//JwtAuth ...
type JwtAuth struct {
	token string
}

// NewJWTAuth ...
func NewJWTAuth(token string) *JwtAuth {
	jwtToken := JwtAuth{token}
	return &jwtToken
}

// GetAuthType ...
func (jwt *JwtAuth) GetAuthType() AuthType {
	return AuthJwt
}

// GetAuthKey ...
func (jwt *JwtAuth) GetAuthKey() (string, error) {
	token := jwt.token
	if token == "" {
		return "", errors.New("token can't be empty")
	}
	token = fmt.Sprintf("Bearer %s", jwt.token)
	return token, nil
}
