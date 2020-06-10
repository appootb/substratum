package rpc

import (
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/errors"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/monitor"
	"github.com/appootb/substratum/recovery"
	"google.golang.org/grpc"

	validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

type ServerOption func(*ServerOptions)

type ServerOptions struct {
	srvOpts []grpc.ServerOption

	unaryChains  []grpc.UnaryServerInterceptor
	streamChains []grpc.StreamServerInterceptor
}

func NewOptions(opts ...ServerOption) *ServerOptions {
	options := &ServerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithServerOption(opts ...grpc.ServerOption) ServerOption {
	return func(options *ServerOptions) {
		options.srvOpts = append(options.srvOpts, opts...)
	}
}

func WithDefaultUnaryInterceptors(fns ...grpc.UnaryServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.unaryChains = append(options.unaryChains,
			recovery.UnaryServerInterceptor(),
			errors.UnaryResponseInterceptor(),
			monitor.UnaryServerInterceptor(),
			metadata.UnaryServerInterceptor(),
			logger.UnaryServerInterceptor(),
			auth.UnaryServerInterceptor(),
			validator.UnaryServerInterceptor(),
		)
		options.unaryChains = append(options.unaryChains, fns...)
	}
}

func WithUnaryInterceptors(fns ...grpc.UnaryServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.unaryChains = append(options.unaryChains, fns...)
	}
}

func WithDefaultStreamInterceptors(fns ...grpc.StreamServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.streamChains = append(options.streamChains,
			recovery.StreamServerInterceptor(),
			errors.StreamServerInterceptor(),
			monitor.StreamServerInterceptor(),
			metadata.StreamServerInterceptor(),
			logger.StreamServerInterceptor(),
			auth.StreamServerInterceptor(),
			validator.StreamServerInterceptor(),
		)
		options.streamChains = append(options.streamChains, fns...)
	}
}

func WithStreamInterceptors(fns ...grpc.StreamServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.streamChains = append(options.streamChains, fns...)
	}
}
