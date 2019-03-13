package netutil

import (
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

func init() {
	SessionFileStore = sessions.NewFilesystemStore("", []byte(os.Getenv("SESSION_KEY")))
	SessionCookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

var (
	SessionFileStore   *sessions.FilesystemStore
	SessionCookieStore *sessions.CookieStore
)

func NewSessionCookieStore() *sessions.CookieStore {
	return SessionCookieStore
}

func NewSessionFileStore() *sessions.FilesystemStore {
	return SessionFileStore
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
