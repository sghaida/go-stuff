package cauth

type AuthType string

const (
	// None no auth provider
	None AuthType = "none"
	// AuthBasic using basic auth
	AuthBasic AuthType = "basic"
	// AuthSts  using DH sts auth
	AuthSts AuthType = "sts"
	// AuthJwt using jwt tokens
	AuthJwt AuthType = "jwt"
	// AuthApiKey using token based auth
	AuthApiKey AuthType = "apikey"
)

type AuthKeyGetter interface {
	GetAuthKey() (string, error)
}

type IAuth interface {
	GetAuthType() AuthType
}

type IBasicAuth interface {
	IAuth
	AuthKeyGetter
}

type IJwtAuth interface {
	IAuth
	AuthKeyGetter
}

type IApiKey interface {
	IAuth
	AuthKeyGetter
}

type INoAuth interface {
	IAuth
}

type NoAuth struct{}

func (n *NoAuth) GetAuthType() AuthType {
	return None
}
