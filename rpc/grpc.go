package rpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

type gRPCServerBuilder struct {
	once sync.Once

	unaryInterceptor  grpc.UnaryServerInterceptor
	streamInterceptor grpc.StreamServerInterceptor
}

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

func New(opts *ServerOptions) *grpc.Server {
	return builder.New(opts)
}

func ChainUnaryServer() grpc.UnaryServerInterceptor {
	return builder.unaryInterceptor
}

func ChainStreamServer() grpc.StreamServerInterceptor {
	return builder.streamInterceptor
}

// chainUnaryServer composes unary interceptors (left-to-right), replacing
// go-grpc-middleware v1 ChainUnaryServer which was removed in v2 in favor of
// grpc.ChainUnaryInterceptor (ServerOption only; gateway still needs one interceptor).
func chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)
	if n == 0 {
		return nil
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		chained := handler
		for i := n - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chained
			chained = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return interceptor(currentCtx, currentReq, info, next)
			}
		}
		return chained(ctx, req)
	}
}

// chainStreamServer composes stream interceptors (left-to-right).
func chainStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	n := len(interceptors)
	if n == 0 {
		return nil
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		chained := handler
		for i := n - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chained
			chained = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
				return interceptor(currentSrv, currentStream, info, next)
			}
		}
		return chained(srv, ss)
	}
}
