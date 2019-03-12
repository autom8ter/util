package netutil

import (
	"encoding/base64"
	"fmt"
	"github.com/autom8ter/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
	"math/rand"
	"net/http"
)

type Router struct {
	addr   string
	router *mux.Router
	chain  *negroni.Negroni
}

func NewRouter(addr string) *Router {
	m := mux.NewRouter()
	n := negroni.Classic()
	return &Router{
		addr:   addr,
		router: m,
		chain:  n,
	}
}
func (r *Router) WithDebug() {
	WithDebug(r.router)
}

func (r *Router) WithPProf() {
	WithPProf(r.router)
}

func (r *Router) WithStatus() {
	WithStatus(r.router)
}

func (r *Router) WithSettings() {
	WithSettings(r.router)
}

func (r *Router) WithStaticViews() {
	WithStaticViews(r.router)
}

func (r *Router) WithRoutes() {
	WithRoutes(r.router)
}

func (r *Router) WithMetrics() {
	WithMetrics(r.router)
}

func (r *Router) BeforeAfter(before, after http.HandlerFunc) {
	r.chain.Use(BeforeNextAfter(before, after))
}

func (r *Router) WithJWT(signingKey string, debug bool, path string, handler http.Handler) {
	WithJWT(signingKey, debug, path, handler, r.router)
}

func (r *Router) Serve() {
	fmt.Printf("starting http server on: %s\n", r.addr)
	r.chain.UseHandler(r.router)
	r.chain.Run(r.addr)
}

func (r *Router) NotImplememntedFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("Not Implemented"))
	}
}

func (r *Router) OnErrorUnauthorized(w http.ResponseWriter, req *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

func (r *Router) OnErrorInternal(w http.ResponseWriter, req *http.Request, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func (r *Router) GenerateJWT(signingKey string, claims map[string]interface{}) (string, error) {
	return util.GenerateJWT(signingKey, claims)
}

func (r *Router) SetResponseHeaders(headers map[string]string, w http.ResponseWriter) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}

func (r *Router) GetHeader(key string, w http.ResponseWriter) string {
	return w.Header().Get(key)
}

func (r *Router) DelHeader(key string, w http.ResponseWriter) {
	w.Header().Del(key)
}

func (r *Router) Do(r2 *http.Request) (*http.Response, error) {
	client := http.DefaultClient
	return client.Do(r2)
}

func (r *Router) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
}

func (r *Router) JSONify(obj interface{}) []byte {
	return util.ToPrettyJson(obj)
}

func (r *Router) RandomTokenString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *Router) RandomToken(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func (r *Router) DerivePassword(counter uint32, password_type, password, user, site string) string {
	return util.DerivePassword(counter, password, password, user, site)
}

func (r *Router) GeneratePrivateKey(typ string) string {
	return util.GeneratePrivateKey(typ)
}

func (r *Router) Render(s string, data interface{}) string {
	return util.Render(s, data)
}

func (r *Router) SetSessionValFunc(cookieStore *sessions.CookieStore, name string, vals map[string]interface{}) http.HandlerFunc {
	return SetSessionValFunc(cookieStore, name, vals)
}

func (r *Router) NewSessionStore(key string) *sessions.CookieStore {
	return NewSessionStore(key)
}

func (r *Router) AddFlashSessionFunc(cookieStore *sessions.CookieStore, name string, val interface{}, vars ...string) http.HandlerFunc {
	return AddFlashSessionFunc(cookieStore, name, val, vars...)
}
