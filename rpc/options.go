package rpc

import (
	"time"

	"github.com/appootb/substratum/v2/auth"
	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/errors"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/monitor"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/recovery"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/tap"

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

func WithDefaultUnaryInterceptors(fns ...grpc.UnaryServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.unaryChains = append(options.unaryChains,
			recovery.UnaryServerInterceptor(),
			errors.UnaryResponseInterceptor(),
			monitor.UnaryServerInterceptor(),
			metadata.UnaryServerInterceptor(),
			auth.UnaryServerInterceptor(),
			logger.UnaryServerInterceptor(),
			validator.UnaryServerInterceptor(),
			client.UnaryServerInterceptor(),
			discovery.UnaryServerInterceptor(),
			storage.UnaryServerInterceptor(),
			queue.UnaryServerInterceptor(),
			task.UnaryServerInterceptor(),
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
			auth.StreamServerInterceptor(),
			logger.StreamServerInterceptor(),
			validator.StreamServerInterceptor(),
			client.StreamServerInterceptor(),
			discovery.StreamServerInterceptor(),
			storage.StreamServerInterceptor(),
			queue.StreamServerInterceptor(),
			task.StreamServerInterceptor(),
		)
		options.streamChains = append(options.streamChains, fns...)
	}
}

func WithStreamInterceptors(fns ...grpc.StreamServerInterceptor) ServerOption {
	return func(options *ServerOptions) {
		options.streamChains = append(options.streamChains, fns...)
	}
}
