package rpc

import (
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var DefaultOptions = NewOptions(WithDefaultUnaryInterceptors(), WithDefaultStreamInterceptors())

type gRPCServerBuilder struct {
	optionsInstalled bool
}

func (b *gRPCServerBuilder) New(opts *ServerOptions) *grpc.Server {
	defer func() {
		b.optionsInstalled = true
	}()

	if b.optionsInstalled {
		return grpc.NewServer()
	}

	if len(opts.unaryChains) > 0 {
		unaryInterceptor := middleware.ChainUnaryServer(opts.unaryChains...)
		opts.srvOpts = append(opts.srvOpts, grpc.UnaryInterceptor(unaryInterceptor))
	}
	if len(opts.streamChains) > 0 {
		streamInterceptor := middleware.ChainStreamServer(opts.streamChains...)
		opts.srvOpts = append(opts.srvOpts, grpc.StreamInterceptor(streamInterceptor))
	}
	return grpc.NewServer(opts.srvOpts...)
}

var builder = &gRPCServerBuilder{}

func New(opts *ServerOptions) *grpc.Server {
	return builder.New(opts)
}
