package util

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

func WithPProf(r *mux.Router) *mux.Router {
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return r
}

// Here we are implementing the NotImplemented handler. Whenever an API endpoint is hit
// we will simply return the message "Not Implemented"
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func ResponseHeaders(headers map[string]string, w http.ResponseWriter) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}

func RequestHeaders(headers map[string]string, r *http.Request) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}

func RequestBasicAuth(userName, password string, r *http.Request) {
	r.SetBasicAuth(userName, password)
}
