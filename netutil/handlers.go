package netutil

import (
	"fmt"
	"github.com/autom8ter/util"
	"github.com/gorilla/mux"
	"net/http"
)

func VarFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/vars")
		w.Write([]byte(util.ToPrettyJsonString(RequestVars(request))))
	}
}

func RouteFunc(r *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			type routeLog struct {
				Name         string
				PathRegExp   string
				PathTemplate string
				HostTemplate string
				Methods      []string
			}
			meth, _ := route.GetMethods()
			host, _ := route.GetHostTemplate()
			pathreg, _ := route.GetPathRegexp()
			pathtemp, _ := route.GetPathTemplate()
			rout := &routeLog{
				Name:         route.GetName(),
				PathRegExp:   pathreg,
				PathTemplate: pathtemp,
				HostTemplate: host,
				Methods:      meth,
			}
			w.Write([]byte(util.ToPrettyJson(rout)))
			fmt.Println("registered handler: ", util.ToPrettyJsonString(rout))
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
