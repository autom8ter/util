package util

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
	"net/http/pprof"
	"os"
)

func WithPProf(r *mux.Router) *mux.Router {
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return r
}

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


func WithLogging(r http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, r)
}

func NotImplememntedFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("Not Implemented"))
	}
}


func OnErrorUnauthorized(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}


func OnErrorInternal(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func WithStatus(r *mux.Router)  {
	r.HandleFunc("/status", func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("API is up and running"))
	})
}
func WithSettings(r *mux.Router)  {
	r.HandleFunc("/settings", func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte(ToPrettyJsonString(viper.AllSettings())))
	})
}
