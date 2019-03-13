package netutil

import (
	"github.com/autom8ter/util/netutil/grpc/api"
	"github.com/autom8ter/util/netutil/grpc/engine"
)

func ServeGrpc(servers ...api.Server) error {
	s := engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			servers...,
		),
	)
	return s.Serve()
}
