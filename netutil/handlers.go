package netutil

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/autom8ter/util"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"net/url"
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
		tokenString, err := util.GenerateJWT(signingKey, claims)
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

func Auth0Login(clientId, clientSecret, redirect, domain, audience string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conf := &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirect,
			Scopes:       []string{"openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + domain + "/authorize",
				TokenURL: "https://" + domain + "/oauth/token",
			},
		}

		if audience == "" {
			audience = "https://" + domain + "/userinfo"
		}

		// Generate random state
		b := make([]byte, 32)
		rand.Read(b)
		state := base64.StdEncoding.EncodeToString(b)

		session, err := SessionFileStore.Get(r, "state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["state"] = state
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		audience := oauth2.SetAuthURLParam("audience", audience)
		url := conf.AuthCodeURL(state, audience)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func Auth0Logout(redirect, clientId, domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Url *url.URL
		Url, err := url.Parse("https://" + domain)

		if err != nil {
			panic("failed to parse domain")
		}

		Url.Path += "/v2/logout"
		parameters := url.Values{}
		parameters.Add("returnTo", redirect)
		parameters.Add("client_id", clientId)
		Url.RawQuery = parameters.Encode()

		http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)
	}
}

func Auth0Callback(sessionName, clientId, clientSecret, redirect, domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conf := &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirect,
			Scopes:       []string{"openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + domain + "/authorize",
				TokenURL: "https://" + domain + "/oauth/token",
			},
		}
		state := r.URL.Query().Get("state")
		session, err := SessionFileStore.Get(r, "state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != session.Values["state"] {
			http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
			return
		}

		code := r.URL.Query().Get("code")

		token, err := conf.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Getting now the userInfo
		client := conf.Client(context.TODO(), token)
		resp, err := client.Get("https://" + domain + "/userinfo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		var profile map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err = SessionFileStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["id_token"] = token.Extra("id_token")
		session.Values["access_token"] = token.AccessToken
		session.Values["profile"] = profile
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to logged in page
		http.Redirect(w, r, "/user", http.StatusSeeOther)

	}
}
