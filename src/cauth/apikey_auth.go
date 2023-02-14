package cauth

import "errors"

type ApiKey struct {
	key string
}

func NewAPIKey(apiKey string) *ApiKey {
	key := ApiKey{apiKey}
	return &key
}

func (k *ApiKey) GetAuthType() AuthType {
	return AuthApiKey
}

func (k *ApiKey) GetAuthKey() (string, error) {
	if k.key == "" {
		return "", errors.New("key should not be empty")
	}
	return k.key, nil
}
