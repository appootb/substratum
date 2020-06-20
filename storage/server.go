package storage

import (
	"context"

	"google.golang.org/grpc"
)

type storageKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for storage manager.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, storageKey{}, DefaultManager), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for storage manager.
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
	return context.WithValue(ctx, storageKey{}, DefaultManager)
}

func ContextStorage(ctx context.Context, component string) Storage {
	if mgr := ctx.Value(storageKey{}); mgr != nil {
		return mgr.(Manager).Get(component)
	}
	return nil
}
