package substratum

import (
	"context"
	"time"

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
		s.ctx = ctx
	}
}

func WithKeepAliveTTL(ttl time.Duration) ServerOption {
	return func(s *Server) {
		s.keepAliveTTL = ttl
	}
}
