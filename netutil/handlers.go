package netutil

import (
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/autom8ter/util"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"strings"
)

func VarFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/vars")
		w.Write([]byte(util.ToPrettyJsonString(RequestVars(request))))
	}
}

func RouteFunc(r *mux.Router) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				type routeLog struct {
					Name         string
					PathRegExp   string
					PathTemplate string
					HostTemplate string
					Methods      []string
				}
				meth, _ := route.GetMethods()
				host, _ := route.GetHostTemplate()
				pathreg, _ := route.GetPathRegexp()
				pathtemp, _ := route.GetPathTemplate()
				rout := &routeLog{
					Name:         route.GetName(),
					PathRegExp:   pathreg,
					PathTemplate: pathtemp,
					HostTemplate: host,
					Methods:      meth,
				}
				w.Write([]byte(util.ToPrettyJson(rout)))
				fmt.Println("registered handler: ", util.ToPrettyJsonString(rout))
				return nil
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
	}
}

func JWTFunc(signingKey string, claims map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Sign the token with our secret */
		tokenString, err := GenerateJWT(signingKey, claims)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		/* Finally, write the token to the browser window */
		w.Write([]byte(tokenString))
	}
}


func Auth0Middleware(clientSecret, domain string, audience []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret = []byte(clientSecret)
		var issuer string
		switch {
		case strings.Contains("https", domain):
			issuer = domain + ".auth0.com/"
		default:
			issuer = "https://" + domain + ".auth0.com/"
		}
		secretProvider := auth0.NewKeyProvider(secret)

		configuration := auth0.NewConfiguration(secretProvider, audience, issuer, jose.HS256)
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
