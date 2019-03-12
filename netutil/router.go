package netutil

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

type Router struct {
	addr string
	router *mux.Router
	chain *negroni.Negroni
}

func NewRouter(addr string) *Router{
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


func(r *Router)  WithStatus() {
	WithStatus(r.router)
}

func (r *Router) WithSettings() {
	WithSettings(r.router)
}

func (r *Router) WithStaticViews() {
	WithStaticViews(r.router)
}


func (r *Router)  WithRoutes() {
	WithRoutes(r.router)
}

func (r *Router)  WithMetrics() {
	WithMetrics(r.router)
}

func (r *Router)  BeforeAfter(before, after http.HandlerFunc) {
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