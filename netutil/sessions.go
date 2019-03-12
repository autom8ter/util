package netutil

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func NewSessionStore(key string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(key))
}

func SetSessionValFunc(cookieStore *sessions.CookieStore, name string, vals map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cookieStore.Get(r, name)
		if err != nil {
			OnErrorUnauthorized(w, r, err.Error())
			return
		}
		for k, v := range vals {
			session.Values[k] = v
		}
		session.Save(r, w)
	}
}

func AddFlashSessionFunc(cookieStore *sessions.CookieStore, name string, val interface{}, vars ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cookieStore.Get(r, name)
		if err != nil {
			OnErrorUnauthorized(w, r, err.Error())
			return
		}
		session.AddFlash(val, vars...)
		session.Save(r, w)
	}
}
