package substratum

import (
	"context"
	"errors"
	"net/http"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/server"
	"github.com/appootb/substratum/storage"
	"google.golang.org/grpc"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	cs map[string]Component
	ms map[permission.VisibleScope]*server.ServeMux
}

func NewServer(opts ...ServerOption) Service {
	srv := &Server{
		cs: make(map[string]Component),
		ms: make(map[permission.VisibleScope]*server.ServeMux),
	}
	for _, opt := range opts {
		opt(srv)
	}
	if srv.ctx == nil {
		srv.ctx, srv.cancel = context.WithCancel(context.Background())
	}
	return srv
}

func (s *Server) HandleFunc(scope permission.VisibleScope, pattern string, handler http.HandlerFunc) {
	if m, ok := s.ms[scope]; ok {
		m.HandleFunc(pattern, handler)
		return
	}
	if scope != permission.VisibleScope_ALL_SCOPES {
		return
	}
	for _, m := range s.ms {
		m.HandleFunc(pattern, handler)
	}
}

func (s *Server) Handle(scope permission.VisibleScope, pattern string, handler http.Handler) {
	if mux, ok := s.ms[scope]; ok {
		mux.Handle(pattern, handler)
		return
	}
	if scope != permission.VisibleScope_ALL_SCOPES {
		return
	}
	for _, mux := range s.ms {
		mux.Handle(pattern, handler)
	}
}

// Get gRPC server of the specified visible scope.
func (s *Server) GetScopedGRPCServer(scope permission.VisibleScope) []*grpc.Server {
	if mux, ok := s.ms[scope]; ok {
		return []*grpc.Server{
			mux.RPCServer(),
		}
	}
	if scope != permission.VisibleScope_ALL_SCOPES {
		return []*grpc.Server{}
	}
	srv := make([]*grpc.Server, 0, len(s.ms))
	for _, mux := range s.ms {
		srv = append(srv, mux.RPCServer())
	}
	return srv
}

// Register gateway handler to the specified visible scope.
func (s *Server) RegisterGateway(scope permission.VisibleScope, handler service.GatewayHandler) error {
	if mux, ok := s.ms[scope]; ok {
		err := handler(s.ctx, mux.GatewayMux(), mux.Client())
		if err != nil {
			return err
		}
	}
	if scope != permission.VisibleScope_ALL_SCOPES {
		return nil
	}
	for _, mux := range s.ms {
		err := handler(s.ctx, mux.GatewayMux(), mux.Client())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) AddMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) error {
	var err error
	if _, ok := s.ms[scope]; ok {
		return errors.New("ServerMux for the specified scope has already been registered")
	}
	s.ms[scope], err = server.NewServeMux(rpcPort, gatewayPort)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Register(comp Component) error {
	s.cs[comp.Name()] = comp

	// Init component.
	if err := comp.Init(s.ctx); err != nil {
		return err
	}
	if err := comp.InitStorage(storage.Default); err != nil {
		return err
	}
	if err := comp.RegisterService(auth.Default, s); err != nil {
		return err
	}
	return nil
}

func (s *Server) Start() {
	for _, mux := range s.ms {
		mux.Serve()
	}
}

func (s *Server) Stop() {
	s.cancel()
}
