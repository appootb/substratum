package task

import (
	"context"

	"google.golang.org/grpc"
)

type taskKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for crontab task.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ContextWithTaskService(ctx), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for crontab task.
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
	return ContextWithTaskService(ctx)
}

func ContextWithTaskService(ctx context.Context) context.Context {
	return context.WithValue(ctx, taskKey{}, impl)
}

func ContextTaskService(ctx context.Context) Task {
	if srv := ctx.Value(taskKey{}); srv != nil {
		return srv.(Task)
	}
	return nil
}
