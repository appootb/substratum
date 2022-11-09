package substratum

import (
	"context"
	"errors"
	"time"

	"github.com/appootb/substratum/v2/auth"
	"github.com/appootb/substratum/v2/configure"
	"github.com/appootb/substratum/v2/discovery"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/plugin"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/rpc"
	"github.com/appootb/substratum/v2/server"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"github.com/appootb/substratum/v2/util/snowflake"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	ctx          context.Context
	keepAliveTTL time.Duration

	components  []Component
	rpcServices map[string][]string
	serveMuxers map[permission.VisibleScope]*server.ServeMux
}

func NewServer(opts ...ServerOption) Service {
	// Register plugin.
	plugin.Register()
	// New server
	srv := &Server{
		ctx:          ictx.Context,
		keepAliveTTL: 3 * time.Second,
		rpcServices:  make(map[string][]string),
		serveMuxers:  make(map[permission.VisibleScope]*server.ServeMux),
	}
	opts = append(opts, WithDefaultClientMux(), WithDefaultServerMux())
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

// Context returns the service context.
func (s *Server) Context() context.Context {
	return s.ctx
}

// UnaryInterceptor returns the unary server interceptor for local gateway handler server.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return rpc.ChainUnaryServer()
}

// StreamInterceptor returns the stream server interceptor for local gateway handler server.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return rpc.ChainStreamServer()
}

// GetGRPCServer returns gRPC server of the specified visible scope.
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

// GetGatewayMux returns gateway mux of the specified visible scope.
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

func (s *Server) AddServeMux(scope permission.VisibleScope, rpcPort, gatewayPort uint16) error {
	var err error
	if _, ok := s.serveMuxers[scope]; ok {
		return errors.New("ServerMux for the specified scope has already been registered")
	}
	metrics := scope == permission.VisibleScope_SERVER
	s.serveMuxers[scope], err = server.NewServeMux(rpcPort, gatewayPort, metrics)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Register(comp Component, rpcs ...string) error {
	name := comp.Name()
	s.components = append(s.components, comp)
	s.rpcServices[name] = rpcs
	storage.Implementor().New(name)

	// Init component.
	if err := comp.Init(configure.Implementor()); err != nil {
		return err
	}
	if err := comp.InitStorage(storage.Implementor().Get(name)); err != nil {
		return err
	}
	if err := comp.RegisterHandler(s.serveMuxers[permission.VisibleScope_CLIENT].HTTPMux(name),
		s.serveMuxers[permission.VisibleScope_SERVER].HTTPMux(name)); err != nil {
		return err
	}
	if err := comp.RegisterService(auth.Implementor(), s); err != nil {
		return err
	}
	return nil
}

func (s *Server) Serve(isolate ...bool) error {
	// Start queue worker and cron tasks.
	for _, comp := range s.components {
		if err := comp.RunQueueWorker(queue.Implementor()); err != nil {
			return err
		}
		if err := comp.ScheduleCronTask(task.Implementor()); err != nil {
			return err
		}
	}

	// Serve muxers.
	for _, mux := range s.serveMuxers {
		mux.Serve()
	}

	// Register node.
	addr := s.serveMuxers[permission.VisibleScope_SERVER].ConnAddr()
	for _, comp := range s.components {
		nodeID, err := discovery.Implementor().Register(comp.Name(), addr,
			discovery.WithIsolate(len(isolate) > 0 && isolate[0]),
			discovery.WithTTL(s.keepAliveTTL),
			discovery.WithServices(s.rpcServices[comp.Name()]))
		if err != nil {
			return err
		}
		// TODO:
		// NodeID is unique in component scope on different node.
		// If multiple components are registered within an unique server,
		// snowflake's node ID might be the same on different nodes.
		snowflake.SetPartitionID(nodeID)
	}

	// Wait for cancellation.
	<-s.ctx.Done()
	return nil
}
