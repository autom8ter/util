package netutil_test

import (
	"fmt"
	"github.com/autom8ter/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ServerTest(handler http.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(util.ToPrettyJsonString(resp))
	}
}

