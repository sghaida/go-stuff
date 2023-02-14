package cauth

import (
	"errors"
	"fmt"
)

type JwtAuth struct {
	token string
}

func NewJWTAuth(token string) *JwtAuth {
	token = fmt.Sprintf("Bearer %s", token)
	jwtToken := JwtAuth{token}
	return &jwtToken
}

func (jwt *JwtAuth) GetAuthType() AuthType {
	return AuthJwt
}

func (jwt *JwtAuth) GetAuthKey() (string, error) {
	if jwt.token == "" {
		return "", errors.New("token can't be empty")
	}
	return jwt.token, nil
}
