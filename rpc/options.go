package rpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/tap"
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

func WithRateLimitingOption(fn tap.ServerInHandle) ServerOption {
	return func(options *ServerOptions) {
		options.srvOpts = append(options.srvOpts, grpc.InTapHandle(fn))
	}
}

func WithDefaultKeepaliveOption() ServerOption {
	return func(options *ServerOptions) {
		options.srvOpts = append(options.srvOpts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     time.Hour,
				MaxConnectionAge:      time.Hour,
				MaxConnectionAgeGrace: time.Hour,
				Time:                  time.Second * 10,
				Timeout:               time.Second * 10,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             time.Second * 10,
				PermitWithoutStream: true,
			}),
		)
	}
}

// WithDefaultUnaryInterceptors installs the framework unary chain merged with
// phase interceptors registered before plugin.Register finishes. Prefer
// RegisterUnary from plugin Init; BuildInterceptorChains runs in plugin.Register.
func WithDefaultUnaryInterceptors() ServerOption {
	return func(options *ServerOptions) {
		options.unaryChains = append(options.unaryChains, UnaryChain()...)
	}
}

// WithUnaryInterceptors appends interceptors after the full default chain.
// Prefer RegisterUnary for plugin-style registration.
func WithUnaryInterceptors(fns ...grpc.UnaryServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.unaryChains = append(options.unaryChains, fns...)
	}
}

// WithDefaultStreamInterceptors installs the framework stream chain merged with
// phase interceptors registered before plugin.Register finishes.
func WithDefaultStreamInterceptors() ServerOption {
	return func(options *ServerOptions) {
		options.streamChains = append(options.streamChains, StreamChain()...)
	}
}

// WithStreamInterceptors appends interceptors after the full default stream chain.
func WithStreamInterceptors(fns ...grpc.StreamServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.streamChains = append(options.streamChains, fns...)
	}
}
