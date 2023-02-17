package cauth

import "errors"

type ApiKey struct {
	key string
}

// NewAPIKey ...
func NewAPIKey(apiKey string) *ApiKey {
	key := ApiKey{apiKey}
	return &key
}

// GetAuthType ...
func (k *ApiKey) GetAuthType() AuthType {
	return AuthApiKey
}

// GetAuthData ...
func (k *ApiKey) GetAuthData() (AuthHeader, error) {
	if k.key == "" {
		return AuthHeader{}, errors.New("key should not be empty")
	}
	return AuthHeader{key: "x-api-key", value: k.key}, nil
}
