package metadata

import (
	"context"
	"sync"

	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type metadataKey struct{}

// holder keeps raw gRPC metadata and parses into *common.Metadata once.
type holder struct {
	raw    metadata.MD
	peerIP string

	once   sync.Once
	parsed *common.Metadata
}

func (h *holder) get() *common.Metadata {
	h.once.Do(func() {
		h.parsed = parseFromMD(h.raw, h.peerIP)
	})
	return h.parsed
}

// ContextWithIncomingMetadata attaches a lazy metadata holder when the context
// carries gRPC incoming metadata. Parsing is deferred until IncomingMetadata.
func ContextWithIncomingMetadata(ctx context.Context) context.Context {
	if ctx.Value(metadataKey{}) != nil {
		return ctx
	}
	raw, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	return context.WithValue(ctx, metadataKey{}, &holder{
		raw:    raw,
		peerIP: peerAddr(ctx),
	})
}

// UnaryServerInterceptor attaches request metadata lazily (no parse on entry).
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ContextWithIncomingMetadata(ctx), req)
	}
}

// StreamServerInterceptor attaches a single lazy holder for the stream.
// Context() reuses that holder so metadata is parsed at most once.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &ctxWrapper{
			ServerStream: stream,
			ctx:          ContextWithIncomingMetadata(stream.Context()),
		})
	}
}

type ctxWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *ctxWrapper) Context() context.Context {
	return s.ctx
}

// IncomingMetadata returns parsed request metadata, parsing lazily on first use
// when the context was prepared by the metadata interceptors.
// Also accepts a pre-parsed *common.Metadata for tests/manual injection.
func IncomingMetadata(ctx context.Context) *common.Metadata {
	v := ctx.Value(metadataKey{})
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case *holder:
		return x.get()
	case *common.Metadata:
		return x
	default:
		return nil
	}
}
