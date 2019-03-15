package util

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpguts"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

//ResponseMiddleware is a function used to modify the response of a reverse proxy
type ResponseMiddleware func(func(response *http.Response) error) func(response *http.Response) error

//RequestMiddleware is a function used to modify the incoming request of a reverse proxy from a client
type RequestMiddleware func(func(req *http.Request)) func(req *http.Request)

//TransportMiddleware is a function used to modify the http RoundTripper that is used by a reverse proxy. The default RoundTripper is initially http.DefaultTransport
type TransportMiddleware func(tripper http.RoundTripper) http.RoundTripper

type ResponseFunc func(*http.Response) error
type ErrorFunc func(http.ResponseWriter, *http.Request, error)
type RequestFunc func(*http.Request)

// Middleware is signature of all http server-side middleware.
type HandlerFunc func(http.Handler) http.Handler
type MiddlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(rw, r, next)
}

// RoundTripperFunc wraps a func to make it into a http.RoundTripper. Similar to http.HandleFunc.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Tripperware is a signature for all http client-side middleware.
type Tripperware func(http.RoundTripper) http.RoundTripper

// WrapClient takes an http.Client and wraps its transport in the chain of tripperwares.
func WrapClient(client *http.Client, wares ...Tripperware) *http.Client {
	if len(wares) == 0 {
		return client
	}

	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	for i := len(wares) - 1; i >= 0; i-- {
		transport = wares[i](transport)
	}

	clone := *client
	clone.Transport = transport
	return &clone
}

func SetResponseHeaders(headers map[string]string, w http.ResponseWriter) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}
func GetHeader(key string, w http.ResponseWriter) string {
	return w.Header().Get(key)
}

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func SetHeaders(headers map[string]string, req *http.Request) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func SetForm(vals map[string]string, req *http.Request) {
	for k, v := range vals {
		req.Form.Set(k, v)
	}
}

func SetBasicAuth(user, password string, req *http.Request) *http.Request {
	req.SetBasicAuth(user, password)
	return req
}

func NewRequest(method, url, user, password string, headers map[string]string, form map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	if user != "" || password != "" {
		req.SetBasicAuth(user, password)
	}
	if headers != nil {
		SetHeaders(headers, req)
	}
	if form != nil {
		SetForm(form, req)
	}
	return req, nil
}

func NewRequestCtx(ctx context.Context, method, url, user, password string, headers map[string]string, form map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	if user != "" || password != "" {
		req.SetBasicAuth(user, password)
	}
	if headers != nil {
		SetHeaders(headers, req)
	}
	if form != nil {
		SetForm(form, req)
	}
	req.WithContext(ctx)
	return req, nil
}
func HTTPErrorHandler(logmsg string, code int) func(rw http.ResponseWriter, req *http.Request, err error) {
	return func(rw http.ResponseWriter, req *http.Request, err error) {
		logrus.Printf("%v\n%v", logmsg, err)
		rw.WriteHeader(code)
	}
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

func ProxyRequestFunc(uRL, method, user, password string, headers map[string]string, form map[string]string) func(req *http.Request) {
	target, err := url.Parse(uRL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = SingleJoiningSlash(target.Path, req.URL.Path)
		req.Method = method
		if user != "" && password != "" {
			req.SetBasicAuth(user, password)
		}
		if form != nil {
			for k, v := range form {
				req.Form.Set(k, v)
			}
		}
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
}

func Ping(endpoint string) error {
	_, err := net.DialTimeout("tcp", endpoint, 250*time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func PingLog(endpoint string, sleep time.Duration) {
	for {
		_, err := net.DialTimeout("tcp", endpoint, 250*time.Millisecond)
		if err != nil {
			log.Println("endpoint unreachable: ", err.Error())
		}
		time.Sleep(sleep)
	}
}

// SignalRunner runs a runner function until an interrupt signal is received, at which point it
// will call stopper.
func SignalRunner(runner, stopper func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		runner()
	}()
	log.Println("hit Ctrl-C to shutdown")
	select {
	case <-signals:
		stopper()
	}
}

// SanitizeApiPrefix forces prefix to be non-empty and end with a slash.
func SanitizeApiPrefix(prefix string) string {
	if len(prefix) == 0 || prefix[len(prefix)-1:] != "/" {
		return prefix + "/"
	}
	return prefix
}

// IsPermanentHTTPHeader checks whether hdr belongs to the list of
// permenant request headers maintained by IANA.
// http://www.iana.org/assignments/message-headers/message-headers.xml
// From https://github.com/grpc-ecosystem/grpc-gateway/blob/7a2a43655ccd9a488d423ea41a3fc723af103eda/runtime/context.go#L157
func IsPermanentHTTPHeader(hdr string) bool {
	switch hdr {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Language",
		"Accept-Ranges",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Cookie",
		"Date",
		"Expect",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Schedule-Tag-Match",
		"If-Unmodified-Since",
		"Max-Forwards",
		"Origin",
		"Pragma",
		"Referer",
		"User-Agent",
		"Via",
		"Warning":
		return true
	}
	return false
}

// IsReserved returns whether the key is reserved by gRPC.
func IsReservedGrpcHeader(key string) bool {
	return strings.HasPrefix(key, "Grpc-")
}

// outgoingHeaderMatcher transforms outgoing metadata into HTTP headers.
// We return any response metadata as is.
func OutgoingHeaderMatcher(metadata string) (string, bool) {
	return metadata, true
}

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func CloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

func RemoveHeaders(h http.Header, headerKey string) {
	if c := h.Get(headerKey); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				h.Del(f)
			}
		}
	}
}

func HeaderContainsToken(header http.Header, headerKey, tokenKey string) bool {
	return httpguts.HeaderValuesContainsToken(header[headerKey], tokenKey)
}

func GetTokenFromHeader(header http.Header, headerKey, tokenKey string) string {
	if !HeaderContainsToken(header, headerKey, tokenKey) {
		logrus.Printf("header key- %s does not contain token key- %s", headerKey, tokenKey)
		return ""
	}
	return header.Get(tokenKey)
}

func ValidHeaderField(s string) bool {
	return httpguts.ValidHeaderFieldName(s)
}

func NewCors(origins, methods, headers []string, creds, options, debug bool, maxAge int) *cors.Cors {
	if methods == nil {
		methods = []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		}
	}
	if origins == nil {
		origins = []string{"*"}
	}
	if headers == nil {
		headers = []string{"*"}
	}
	opts := cors.Options{
		AllowedOrigins:     origins,
		AllowedMethods:     methods,
		AllowedHeaders:     headers,
		MaxAge:             maxAge,
		AllowCredentials:   creds,
		OptionsPassthrough: options,
		Debug:              debug,
	}
	return cors.New(opts)
}

func ExecHandler(ctx context.Context, name, dir string, args ...string) http.HandlerFunc {
	type Command struct {
		Name   string   `json:"name"`
		Dir    string   `json:"dir"`
		Args   []string `json:"args"`
		Output []byte   `json:"output"`
	}
	var cmd = &Command{
		Name: name,
		Dir:  dir,
		Args: args,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		bits, err := Exec(ctx, name, dir, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cmd.Output = bits
		w.Write(ToPrettyJson(cmd))
	}
}

func NewSessionCookieStore(key string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(key))
}

type CorsConfig struct {
	Origins, Methods, Headers []string
	Creds, Options, Debug     bool
	MaxAge                    int
}

