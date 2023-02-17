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

// GetAuthData ...
func (jwt *JwtAuth) GetAuthData() (AuthHeader, error) {
	token := jwt.token
	if token == "" {
		return AuthHeader{}, errors.New("token can't be empty")
	}
	token = fmt.Sprintf("Bearer %s", jwt.token)

	return AuthHeader{key: "Authorization", value: token}, nil
}
