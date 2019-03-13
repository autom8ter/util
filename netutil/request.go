package netutil

import (
	"github.com/autom8ter/util"
	"github.com/autom8ter/util/fsutil"
	"github.com/gorilla/mux"
	"net/http"
)

func RequestVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func FileToRequestBody(file string, req *http.Request) (*http.Request, error) {
	bits, err := fsutil.NewFs().ReadFile(file)
	if err != nil {
		return req, err
	}
	bits = util.ToPrettyJson(bits)
	_, err = req.Body.Read(bits)
	if err != nil {
		return req, err
	}
	return req, nil
}

func SetHeaders(headers map[string]string, req *http.Request) *http.Request {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func SetForm(vals map[string]string, req *http.Request) *http.Request {
	for k, v := range vals {
		req.Form.Set(k, v)
	}
	return req
}

func SetBasicAuth(user, password string, req *http.Request) *http.Request {
	req.SetBasicAuth(user, password)
	return req
}
