package netutil

import (
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/auth0/go-jwt-middleware"
	"github.com/autom8ter/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"os"
	"strings"
	"time"
)

type JWTRouter struct {
	*mux.Router
	SignKey       	string
	Claims        	map[string]interface{}
	ClientID   		string
	ClientSecret 	string
	Audience 		[]string
	Domain        	string
	Callback      	string
}

func NewJWTRouter() *JWTRouter {
	claims := make(map[string]interface{})
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["admin"] = true
	claims["name"] = os.Getenv("USER")
	viper.SetDefault("callback", "localhost:8080/callback")
	viper.SetDefault("signkey", "mysupersecretsigningkey")
	viper.SetDefault("claims", claims)
	return &JWTRouter{
		Router:        mux.NewRouter(),
		SignKey:       viper.GetString("signkey"),
		Claims:        viper.GetStringMap("claims"),
		Audience: viper.GetStringSlice("audience"),
		ClientSecret:   viper.GetString("client-secret"),
		ClientID: viper.GetString("client-id"),
		Domain:        viper.GetString("domain"),
		Callback:      viper.GetString("callback"),
	}

}

func (j *JWTRouter) JWTTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Sign the token with our secret */
		tokenString, err := j.GenerateJWT()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
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
		Debug:               true,
		SigningMethod:       jwt.SigningMethodHS256,
	})
}

func (j *JWTRouter) WithJWT(handler http.Handler) http.Handler {
	return j.JWTMiddleware().Handler(handler)
}

func (j *JWTRouter) Auth0Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret = []byte(j.ClientSecret)
		var issuer string
		switch {
		case strings.Contains("https", j.Domain):
			issuer = j.Domain + ".auth0.com/"
		default:
			issuer = "https://" + j.Domain + ".auth0.com/"
		}
		secretProvider := auth0.NewKeyProvider(secret)

		configuration := auth0.NewConfiguration(secretProvider, j.Audience, issuer, jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Unauthorized \n%s\n%s", err.Error(), util.ToPrettyJsonString(r.Header))))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}