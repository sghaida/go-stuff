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

var NoAuth = noAuth{}

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

// ISomeAuth has some auth implementation
type ISomeAuth interface {
	IAuth
	AuthDataGetter
}

// INoAuth has no auth implementation
type INoAuth interface {
	IAuth
}

type noAuth struct{}

func (n noAuth) GetAuthType() AuthType {
	return None
}
