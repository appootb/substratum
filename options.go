package substratum

import (
	"context"

	"github.com/appootb/protobuf/go/permission"
)

type ServerOption func(*Server)

func WithMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) ServerOption {
	return func(s *Server) {
		err := s.AddMux(scope, rpcPort, gatewayPort)
		if err != nil {
			panic(err)
		}
	}
}

func WithContext(ctx context.Context) ServerOption {
	return func(s *Server) {
		s.ctx, s.cancel = context.WithCancel(ctx)
	}
}
