package cauth

type AuthType string

const (
	// None no auth provider
	None AuthType = "none"
	// AuthBasic using basic auth
	AuthBasic AuthType = "basic"

	// AuthJwt using jwt tokens
	AuthJwt AuthType = "jwt"
	// AuthApiKey using token based auth
	AuthApiKey AuthType = "apikey"
)

type AuthHeader struct {
	key   string
	value string
}

func NewAuthHeader(key string, value string) AuthHeader {
	return AuthHeader{
		key:   key,
		value: value,
	}
}

func (a AuthHeader) GetAuthKeyValue() (string, string) {
	return a.key, a.value
}

type AuthDataGetter interface {
	GetAuthData() (AuthHeader, error)
}

type IAuth interface {
	GetAuthType() AuthType
}

type IBasicAuth interface {
	IAuth
	AuthDataGetter
}

type IJwtAuth interface {
	IAuth
	AuthDataGetter
}

type IApiKey interface {
	IAuth
	AuthDataGetter
}

type INoAuth interface {
	IAuth
}

type NoAuth struct{}

func (n *NoAuth) GetAuthType() AuthType {
	return None
}
