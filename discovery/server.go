package discovery

import (
	"context"

	"google.golang.org/grpc"
)

type discoveryKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for service discovery.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ContextWithDiscovery(ctx), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for service discovery.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &ctxWrapper{stream}
		return handler(srv, wrapper)
	}
}

type ctxWrapper struct {
	grpc.ServerStream
}

func (s *ctxWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()
	return ContextWithDiscovery(ctx)
}

func ContextWithDiscovery(ctx context.Context) context.Context {
	return context.WithValue(ctx, discoveryKey{}, impl)
}

func ContextDiscovery(ctx context.Context) Discovery {
	if mgr := ctx.Value(discoveryKey{}); mgr != nil {
		return mgr.(Discovery)
	}
	return nil
}
