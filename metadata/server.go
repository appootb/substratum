package metadata

import (
	"context"

	"github.com/appootb/protobuf/go/common"
	"google.golang.org/grpc"
)

type metadataKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for parsing request metadata.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := ParseIncomingMetadata(ctx)
		return handler(context.WithValue(ctx, metadataKey{}, md), req)
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
	md := ParseIncomingMetadata(ctx)
	return context.WithValue(ctx, metadataKey{}, md)
}

func RequestMetadata(ctx context.Context) *common.Metadata {
	if md := ctx.Value(metadataKey{}); md != nil {
		return md.(*common.Metadata)
	}
	return ParseIncomingMetadata(ctx)
}
