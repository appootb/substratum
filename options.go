package substratum

import (
	"context"
	"time"

	"github.com/appootb/protobuf/go/permission"
)

const (
	DefaultClientRpcPort     = 8080
	DefaultClientGatewayPort = 8081
	DefaultServerRpcPort     = 8088
	DefaultServerGatewayPort = 8089
)

type ServerOption func(*Server)

func WithDefaultClientMux() ServerOption {
	return WithMux(permission.VisibleScope_CLIENT, DefaultClientRpcPort, DefaultClientGatewayPort)
}

func WithDefaultServerMux() ServerOption {
	return WithMux(permission.VisibleScope_SERVER, DefaultServerRpcPort, DefaultServerGatewayPort)
}

func WithMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) ServerOption {
	return func(s *Server) {
		if _, ok := s.serveMuxers[scope]; ok {
			return
		}
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
