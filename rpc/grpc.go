package rpc

import (
	"sync"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
			b.unaryInterceptor = middleware.ChainUnaryServer(opts.unaryChains...)
		}
		if len(opts.streamChains) > 0 {
			b.streamInterceptor = middleware.ChainStreamServer(opts.streamChains...)
		}
	})

	opts.srvOpts = append(opts.srvOpts, grpc.UnaryInterceptor(b.unaryInterceptor))
	opts.srvOpts = append(opts.srvOpts, grpc.StreamInterceptor(b.streamInterceptor))
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
