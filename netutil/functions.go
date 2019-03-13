package netutil

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"strings"
	"sync"
)

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

func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// PassedHeaderDeciderFunc returns true if given header should be passed to gRPC server metadata.
type PassedHeaderDeciderFunc func(string) bool

func CreatePassingHeaderMiddleware(decide PassedHeaderDeciderFunc) HandlerFunc {
	return func(next http.Handler) http.Handler {
		cache := new(sync.Map)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newHeader := make(http.Header, 2*len(r.Header))

			for k := range r.Header {
				v := r.Header.Get(k)
				if newKey, ok := cache.Load(k); ok {
					newHeader.Set(newKey.(string), v)
				} else if decide(k) {
					newKey := runtime.MetadataHeaderPrefix + k
					cache.Store(k, newKey)
					newHeader.Set(newKey, v)
				}
				newHeader.Set(k, v)
			}

			r.Header = newHeader

			next.ServeHTTP(w, r)
		})
	}
}
