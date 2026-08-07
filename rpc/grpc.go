package rpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// gRPCServerBuilder builds a process-wide gRPC server interceptor chain once.
type gRPCServerBuilder struct {
	once sync.Once

	unaryInterceptor  grpc.UnaryServerInterceptor
	streamInterceptor grpc.StreamServerInterceptor
}

// New creates a *grpc.Server with the composed unary/stream interceptors from opts.
func (b *gRPCServerBuilder) New(opts *ServerOptions) *grpc.Server {
	b.once.Do(func() {
		if len(opts.unaryChains) > 0 {
			b.unaryInterceptor = chainUnaryServer(opts.unaryChains...)
		}
		if len(opts.streamChains) > 0 {
			b.streamInterceptor = chainStreamServer(opts.streamChains...)
		}
	})

	if b.unaryInterceptor != nil {
		opts.srvOpts = append(opts.srvOpts, grpc.UnaryInterceptor(b.unaryInterceptor))
	}
	if b.streamInterceptor != nil {
		opts.srvOpts = append(opts.srvOpts, grpc.StreamInterceptor(b.streamInterceptor))
	}
	return grpc.NewServer(opts.srvOpts...)
}

var builder = &gRPCServerBuilder{}

// New creates a gRPC server using the package-level builder and opts.
func New(opts *ServerOptions) *grpc.Server {
	return builder.New(opts)
}

// ChainUnaryServer returns the composed unary interceptor after the first New call.
// Used by local gateway handlers that need the same chain as the gRPC server.
func ChainUnaryServer() grpc.UnaryServerInterceptor {
	return builder.unaryInterceptor
}

// ChainStreamServer returns the composed stream interceptor after the first New call.
func ChainStreamServer() grpc.StreamServerInterceptor {
	return builder.streamInterceptor
}

// chainUnaryServer composes unary interceptors (left-to-right), replacing
// go-grpc-middleware v1 ChainUnaryServer which was removed in v2 in favor of
// grpc.ChainUnaryInterceptor (ServerOption only; gateway still needs one interceptor).
//
// Uses a single per-request continuation closure instead of allocating one
// closure per interceptor on every RPC. Each interceptor must call next
// synchronously and exactly once per invocation; do not retain or reuse next.
func chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)
	switch n {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		i := 0
		var next grpc.UnaryHandler
		next = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
			if i == n {
				return handler(currentCtx, currentReq)
			}
			curr := interceptors[i]
			i++
			return curr(currentCtx, currentReq, info, next)
		}
		return next(ctx, req)
	}
}

// chainStreamServer composes stream interceptors (left-to-right).
// Same single-continuation strategy and next-call contract as chainUnaryServer.
func chainStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	n := len(interceptors)
	switch n {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		i := 0
		var next grpc.StreamHandler
		next = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
			if i == n {
				return handler(currentSrv, currentStream)
			}
			curr := interceptors[i]
			i++
			return curr(currentSrv, currentStream, info, next)
		}
		return next(srv, ss)
	}
}
