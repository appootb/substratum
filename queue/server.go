package queue

import (
	"context"

	"google.golang.org/grpc"
)

type queueKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for message queue service.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, queueKey{}, impl), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for message queue service.
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
	return context.WithValue(ctx, queueKey{}, impl)
}

func ContextQueueService(ctx context.Context) Service {
	if srv := ctx.Value(queueKey{}); srv != nil {
		return srv.(Service)
	}
	return nil
}
