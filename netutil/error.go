package netutil

import "net/http"

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
