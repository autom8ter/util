package util

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
)

func WithPProf(r *mux.Router) *mux.Router {
	fmt.Println("registered handler: ", "/debug/pprof/")
	r.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	fmt.Println("registered handler: ", "/debug/pprof/cmdline")
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	fmt.Println("registered handler: ", "/debug/pprof/profile")
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	fmt.Println("registered handler: ", "/debug/pprof/symbol")
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	fmt.Println("registered handler: ", "/debug/pprof/trace")
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

func WithStatus(r *mux.Router) {
	r.HandleFunc("/status", func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/status")
		w.Write([]byte("API is up and running"))
	})
}
func WithSettings(r *mux.Router) {
	r.HandleFunc("/settings", func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/settings")
		w.Write([]byte(ToPrettyJsonString(viper.AllSettings())))
	})
}

func WithVars(r *mux.Router) {
	r.HandleFunc("/vars", func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/vars")
		w.Write([]byte(ToPrettyJsonString(RequestVars(request))))
	})
}

func WithStaticViews(r *mux.Router) {
	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	fmt.Println("registered file server handler: ", "./views/")
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	fmt.Println("registered file server handler: ", "./static/")
}

func RequestVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func LogRoutes(r *mux.Router) {
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		meth, _ := route.GetMethods()
		rout := &routeLog{
			Name:    route.GetName(),
			Methods: meth,
		}
		fmt.Println("Registered Handler: ", ToPrettyJsonString(rout))
		return nil
	}); err != nil {
		log.Fatal(err.Error())
	}
}

type routeLog struct {
	Name    string
	Methods []string
}
