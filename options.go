package substratum

import (
	"context"
	"time"

	"github.com/appootb/protobuf/go/permission"
)

const (
	DefaultRpcPort          = 8080
	DefaultGatewayPort      = 8081
	DefaultInnerRpcPort     = 8088
	DefaultInnerGatewayPort = 8089
)

type ServerOption func(*Server)

func WithDefaultMux() ServerOption {
	return WithMux(permission.VisibleScope_DEFAULT_SCOPE, DefaultRpcPort, DefaultGatewayPort)
}

func WithDefaultInnerMux() ServerOption {
	return WithMux(permission.VisibleScope_INNER_SCOPE, DefaultInnerRpcPort, DefaultInnerGatewayPort)
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
