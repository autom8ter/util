package netutil

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RequestHeaders(headers map[string]string, r *http.Request) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}

func RequestBasicAuth(userName, password string, r *http.Request) {
	r.SetBasicAuth(userName, password)
}


func RequestVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}


