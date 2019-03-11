package netutil

import (
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"strings"
)

type JWTRouter struct {
	*mux.Router
	SignKey       string
	Claims        map[string]interface{}
	OAuthSecret   string
	OAuthAudience []string
	Domain        string
	Callback      string
}

func NewJWTRouter() *JWTRouter {
	return &JWTRouter{
		Router:        mux.NewRouter(),
		SignKey:       viper.GetString("jwt.signing-key"),
		Claims:        viper.GetStringMap("jwt.claims"),
		OAuthAudience: viper.GetStringSlice("oauth.audience"),
		OAuthSecret:   viper.GetString("oauth.secret"),
		Domain:        viper.GetString("oauth.domain"),
		Callback:      viper.GetString("callback"),
	}

}

func (j *JWTRouter) JWTTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Sign the token with our secret */
		tokenString, _ := j.GenerateJWT()

		/* Finally, write the token to the browser window */
		w.Write([]byte(tokenString))
	}
}

func (j *JWTRouter) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	clms := token.Claims.(jwt.MapClaims)

	for k, v := range j.Claims {
		clms[k] = v
	}

	tokenString, err := token.SignedString(j.SignKey)

	if err != nil {
		return "", fmt.Errorf("Something Went Wrong: %s", err.Error())
	}

	return tokenString, nil
}

func (j *JWTRouter) JWTMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return j.SignKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func (j *JWTRouter) WithJWT(path string, method []string, handler http.HandlerFunc) {
	j.Handle(path, j.JWTMiddleware().Handler(handler)).Methods(method...)
}

func (j *JWTRouter) Auth0Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret = []byte(j.OAuthSecret)
		var issuer string
		switch {
		case strings.Contains("https", j.Domain):
			issuer = j.Domain + ".auth0.com/"
		default:
			issuer = "https://" + j.Domain + ".auth0.com/"
		}
		secretProvider := auth0.NewKeyProvider(secret)

		configuration := auth0.NewConfiguration(secretProvider, j.OAuthAudience, issuer, jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
