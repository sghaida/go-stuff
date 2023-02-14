package cauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiKey_GetAuthType(t *testing.T) {
	apikey := NewAPIKey("")
	assert.Equal(t, apikey.GetAuthType(), AuthApiKey)
}

func TestApiKey_GetAuthKey(t *testing.T) {

	tt := []struct {
		name         string
		key          string
		expectsError bool
	}{
		{
			name:         "get some api key successfully",
			key:          "some-api-key",
			expectsError: false,
		}, {
			name:         "no api key provided",
			key:          "",
			expectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			apikey := NewAPIKey(tc.key)
			key, err := apikey.GetAuthKey()
			if tc.expectsError {
				assert.Error(t, err)
				assert.Equal(t, key, "")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.key, key)
		})
	}
}
