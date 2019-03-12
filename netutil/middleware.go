package netutil

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)
type MiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(rw, r, next)
}

func BeforeNextAfter(before, after http.HandlerFunc) MiddlewareFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if before != nil {
			before(rw, r)
		}
		if next != nil {
			next(rw, r)
		}
		if after != nil {
			after(rw, r)
		}
	}
}

func JWTMiddleware(singingKey string, debug bool) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return singingKey, nil
		},
		Debug:         debug,
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func WithJWT(signingKey string, debug bool, path string, handler http.Handler, r *mux.Router) {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		},
		Debug:         debug,
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("code"),
			jwtmiddleware.FromParameter("auth-code")),
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.Handle(path, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler),
	))
}



