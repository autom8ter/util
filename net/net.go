package net

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"time"
	"github.com/gorilla/mux"
)

func ClientAPI(req *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(body))
	}
}

func ClientAPIWithJWT(req *http.Request, signKey, clientName string, exp *time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		validToken, err := GenerateJWT(signKey, clientName, exp)
		if err != nil {
			fmt.Println("Failed to generate token")
		}
		req.Header.Set("Token", validToken)
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(body))
	}
}

func GenerateJWT(signKey, client string, exp *time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = client
	claims["exp"] = exp

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func JWTAuthorized(signKey string, endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return signKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func WithPProf(r *mux.Router) *mux.Router {
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return r
}
