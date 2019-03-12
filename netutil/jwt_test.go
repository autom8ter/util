package netutil_test

import (
	"fmt"
	"github.com/autom8ter/util/netutil"
	"errors"
	"testing"
)

var jwtRouter = netutil.NewJWTRouter()

func TestJWTRouter_GenerateJWT(t *testing.T) {
	s, err := jwtRouter.GenerateJWT()
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal(err.Error())
	}
	if s == "" {
		t.Fatal(errors.New("empty token"))
	}
	fmt.Println(s)
}
