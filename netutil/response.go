package netutil

import "net/http"

func SetResponseHeaders(headers map[string]string, w http.ResponseWriter) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}
func GetHeader(key string, w http.ResponseWriter) string {
	return w.Header().Get(key)
}
