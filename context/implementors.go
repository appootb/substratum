package context

import (
	"context"

	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"google.golang.org/grpc"
)

// WithImplementors puts the process-wide client, discovery, storage, queue, and
// task implementors on ctx so handlers can resolve them (ContextConnPool,
// ContextStorage, …). Prefer this over stacking their individual interceptors.
//
// Unlike WithServerContext, this does not set component name or logger helper;
// those are handled by auth / logger interceptors on the gRPC chain.
func WithImplementors(ctx context.Context) context.Context {
	return client.ContextWithConnPool(
		discovery.ContextWithDiscovery(
			queue.ContextWithQueueService(
				storage.ContextWithStorage(
					task.ContextWithTaskService(ctx)))))
}

// UnaryServerInterceptor attaches implementors for unary RPCs.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(WithImplementors(ctx), req)
	}
}

// StreamServerInterceptor attaches implementors for streaming RPCs.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &implementorStream{
			ServerStream: stream,
			ctx:          WithImplementors(stream.Context()),
		})
	}
}

type implementorStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *implementorStream) Context() context.Context {
	return s.ctx
}
