package cauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicAuth_GetAuthType(t *testing.T) {
	basicAuth := NewBasicAuth("", "")
	assert.Equal(t, basicAuth.GetAuthType(), AuthBasic)
}

func TestBasicAuth_GetAuthKey(t *testing.T) {

	tt := []struct {
		name         string
		username     string
		password     string
		expectsError bool
	}{
		{
			name:         "get some basic auth token successfully",
			username:     "some-user",
			password:     "some-password",
			expectsError: false,
		}, {
			name:         "missing username",
			password:     "some-password",
			expectsError: true,
		}, {
			name:         "missing password",
			username:     "some-username",
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			basicAuth := NewBasicAuth(tc.username, tc.password)
			key, err := basicAuth.GetAuthData()
			if tc.expectsError {
				assert.Error(t, err)
				assert.Equal(t, key, "")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.username, basicAuth.username)
			assert.Equal(t, tc.password, basicAuth.password)
		})
	}
}
