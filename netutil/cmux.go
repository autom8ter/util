package netutil

import (
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
	"net"
)

// MuxServer wraps a connection multiplexer and a listener.
type MuxServer struct {
	mux cmux.CMux
	lis net.Listener
}

// Serve implements Server.Serve
func (s *MuxServer) Serve(net.Listener) error {
	grpclog.Info("mux is starting %s", s.lis.Addr())

	err := s.mux.Serve()

	grpclog.Infof("mux is closed: %v", err)

	return errors.Wrap(err, "failed to serve cmux server")
}

// Shutdown implements Server.Shutdown
func (s *MuxServer) Shutdown() {
	err := s.lis.Close()
	if err != nil {
		grpclog.Errorf("failed to close cmux's listener: %v", err)
	}
}
