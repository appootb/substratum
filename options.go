package substratum

import (
	"context"
	"time"

	"github.com/appootb/substratum/v2/proto/go/permission"
)

const (
	DefaultClientRpcPort     = 8080
	DefaultClientGatewayPort = 8081
	DefaultServerRpcPort     = 8088
	DefaultServerGatewayPort = 8089
)

type ServerOption func(*Server)

func WithDefaultClientMux() ServerOption {
	return WithServeMux(permission.VisibleScope_CLIENT, DefaultClientRpcPort, DefaultClientGatewayPort)
}

func WithDefaultServerMux() ServerOption {
	return WithServeMux(permission.VisibleScope_SERVER, DefaultServerRpcPort, DefaultServerGatewayPort)
}

func WithServeMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) ServerOption {
	return func(s *Server) {
		if _, ok := s.serveMuxers[scope]; ok {
			return
		}
		err := s.AddServeMux(scope, rpcPort, gatewayPort)
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
