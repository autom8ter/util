package gokit

import (
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/url"
)

type ValidateJWT func() error

func (v ValidateJWT) Valid() error {
	return v()
}

func ParseJWT(endpoint endpoint.Endpoint, signKey string) endpoint.Endpoint {
	var kf = func(token *stdjwt.Token) (interface{}, error) { return []byte(signKey), nil }

	endpoint = jwt.NewParser(kf, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(endpoint)
	return endpoint
}

func ClientSignJWT(v ValidateJWT, header, signingString, method string, urrl *url.URL, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint {
	return jwt.NewSigner(header, []byte(signingString), stdjwt.SigningMethodHS256, v)(httptransport.NewClient(method, urrl, enc, dec, httptransport.ClientBefore(jwt.ContextToHTTP())).Endpoint())
}