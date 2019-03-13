package api

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type RegisterServerFunc func(*grpc.Server)
type RegisterHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type Driver struct {
	RegisterServerFunc  func(*grpc.Server)
	RegisterHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

func (d *Driver) RegisterWithServer(s *grpc.Server) {
	d.RegisterServerFunc(s)
}

func (d *Driver) RegisterWithHandler(ctx context.Context, m *runtime.ServeMux, cc *grpc.ClientConn) error {
	return d.RegisterHandlerFunc(ctx, m, cc)
}

// Server is an interface for representing gRPC server implementations.
type Server interface {
	RegisterWithServer(*grpc.Server)
	RegisterWithHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

type Interface interface {
	Serve(l net.Listener) error
	Shutdown()
}
