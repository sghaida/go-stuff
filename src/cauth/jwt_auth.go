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
	token = fmt.Sprintf("Bearer %s", token)
	jwtToken := JwtAuth{token}
	return &jwtToken
}

// GetAuthType ...
func (jwt *JwtAuth) GetAuthType() AuthType {
	return AuthJwt
}

// GetAuthKey ...
func (jwt *JwtAuth) GetAuthKey() (string, error) {
	if jwt.token == "" {
		return "", errors.New("token can't be empty")
	}
	return jwt.token, nil
}
