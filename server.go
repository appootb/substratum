package substratum

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/resolver"
	"github.com/appootb/substratum/rpc"
	"github.com/appootb/substratum/server"
	"github.com/appootb/substratum/storage"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	ctx          context.Context
	keepAliveTTL time.Duration

	components  map[string]Component
	serveMuxers map[permission.VisibleScope]*server.ServeMux
}

func NewServer(opts ...ServerOption) Service {
	srv := &Server{
		ctx:          context.Background(),
		keepAliveTTL: 3 * time.Second,
		components:   make(map[string]Component),
		serveMuxers:  make(map[permission.VisibleScope]*server.ServeMux),
	}
	opts = append(opts, WithDefaultClientMux(), WithDefaultServerMux())
	for _, opt := range opts {
		opt(srv)
	}
	// Register gRPC resolver.
	resolver.Register()
	return srv
}

func (s *Server) HandleFunc(scope permission.VisibleScope, pattern string, handler http.HandlerFunc) {
	if m, ok := s.serveMuxers[scope]; ok {
		m.HandleFunc(pattern, handler)
		return
	}
	if scope != permission.VisibleScope_ALL {
		return
	}
	for _, m := range s.serveMuxers {
		m.HandleFunc(pattern, handler)
	}
}

func (s *Server) Handle(scope permission.VisibleScope, pattern string, handler http.Handler) {
	if mux, ok := s.serveMuxers[scope]; ok {
		mux.Handle(pattern, handler)
		return
	}
	if scope != permission.VisibleScope_ALL {
		return
	}
	for _, mux := range s.serveMuxers {
		mux.Handle(pattern, handler)
	}
}

// Return server context.
func (s *Server) Context() context.Context {
	return s.ctx
}

// Return the unary server interceptor for local gateway handler server.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return rpc.ChainUnaryServer()
}

// Return the stream server interceptor for local gateway handler server.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return rpc.ChainStreamServer()
}

// Get gRPC server of the specified visible scope.
func (s *Server) GetGRPCServer(scope permission.VisibleScope) []*grpc.Server {
	if mux, ok := s.serveMuxers[scope]; ok {
		return []*grpc.Server{
			mux.RPCServer(),
		}
	}
	if scope != permission.VisibleScope_ALL {
		return []*grpc.Server{}
	}
	srv := make([]*grpc.Server, 0, len(s.serveMuxers))
	for _, mux := range s.serveMuxers {
		srv = append(srv, mux.RPCServer())
	}
	return srv
}

// Get gateway mux of the specified visible scope.
func (s *Server) GetGatewayMux(scope permission.VisibleScope) []*runtime.ServeMux {
	if mux, ok := s.serveMuxers[scope]; ok {
		return []*runtime.ServeMux{
			mux.GatewayMux(),
		}
	}
	if scope != permission.VisibleScope_ALL {
		return []*runtime.ServeMux{}
	}
	srv := make([]*runtime.ServeMux, 0, len(s.serveMuxers))
	for _, mux := range s.serveMuxers {
		srv = append(srv, mux.GatewayMux())
	}
	return srv
}

func (s *Server) AddMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) error {
	var err error
	if _, ok := s.serveMuxers[scope]; ok {
		return errors.New("ServerMux for the specified scope has already been registered")
	}
	s.serveMuxers[scope], err = server.NewServeMux(rpcPort, gatewayPort)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Register(comp Component) error {
	name := comp.Name()
	s.components[name] = comp
	storage.DefaultManager.New(name)

	// Init component.
	if err := comp.Init(discovery.DefaultConfig); err != nil {
		return err
	}
	if err := comp.InitStorage(storage.DefaultManager.Get(name)); err != nil {
		return err
	}
	if err := comp.RegisterService(auth.Default, s); err != nil {
		return err
	}
	return nil
}

func (s *Server) Serve() error {
	for _, mux := range s.serveMuxers {
		mux.Serve()
	}
	// Register node.
	addr := s.serveMuxers[permission.VisibleScope_SERVER].ConnAddr()
	for name := range s.components {
		err := discovery.DefaultService.RegisterNode(name, addr, s.keepAliveTTL)
		if err != nil {
			return err
		}
	}
	// Wait for cancellation
	<-s.ctx.Done()
	return nil
}
