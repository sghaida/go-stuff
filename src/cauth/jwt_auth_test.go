package cauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtAuth_GetAuthType(t *testing.T) {
	jwtAuth := NewJWTAuth("")
	assert.Equal(t, jwtAuth.GetAuthType(), AuthJwt)
}

func TestJwtAuth_GetAuthKey(t *testing.T) {

	tt := []struct {
		name         string
		jwt          string
		expectsError bool
	}{
		{
			name:         "get some basic auth token successfully",
			jwt:          "some jwt token",
			expectsError: false,
		}, {
			name:         "missing jwt",
			jwt:          "",
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			jwtAuth := NewJWTAuth(tc.jwt)
			key, err := jwtAuth.GetAuthKey()
			if tc.expectsError {
				assert.Error(t, err)
				assert.Equal(t, key, "")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "Bearer "+tc.jwt, key)
		})
	}
}
