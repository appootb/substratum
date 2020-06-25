package client

import (
	"context"

	"google.golang.org/grpc"
)

type connPoolKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for gRPC client connection pool.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, connPoolKey{}, DefaultConnPool), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for parsing request metadata.
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
	return context.WithValue(ctx, connPoolKey{}, DefaultConnPool)
}

func ContextConnPool(ctx context.Context) ConnPool {
	if mgr := ctx.Value(connPoolKey{}); mgr != nil {
		return mgr.(ConnPool)
	}
	return nil
}
