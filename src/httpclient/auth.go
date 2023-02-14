package httpclient

import (
	"fmt"
	"github.com/sghaida/go-stuff/src/cauth"
)

func (c *Caller) addBasicAuth(authData cauth.IBasicAuth) (*Caller, error) {
	authKey, err := authData.GetAuthKey()
	if err != nil {
		return c, err
	}
	c.authHeader["Authorization"] = fmt.Sprintf("Basic %s", authKey)
	return c, nil
}

func (c *Caller) addJwtToken(authData cauth.IJwtAuth) (*Caller, error) {
	token, err := authData.GetAuthKey()
	if err != nil {
		return c, err
	}
	c.authHeader["Authorization"] = token
	return c, nil
}

func (c *Caller) addApiKey(authData cauth.IApiKey) (*Caller, error) {
	apikey, err := authData.GetAuthKey()
	if err != nil {
		return c, err
	}
	c.authHeader["x-api-key"] = apikey
	return c, nil
}
