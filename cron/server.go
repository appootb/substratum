package cron

import (
	"context"

	"google.golang.org/grpc"
)

type cronKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for crontab service.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, cronKey{}, impl), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for crontab service.
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
	return context.WithValue(ctx, cronKey{}, impl)
}

func ContextCronService(ctx context.Context) Cron {
	if srv := ctx.Value(cronKey{}); srv != nil {
		return srv.(Cron)
	}
	return nil
}
