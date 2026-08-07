package context

import (
	"context"
	"testing"

	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"google.golang.org/grpc"
)

func TestWithImplementorsPreservesValues(t *testing.T) {
	base := context.Background()
	ctx := WithImplementors(base)

	// Each getter should find its key even when implementor is nil.
	_ = client.ContextConnPool(ctx)
	_ = discovery.ContextDiscovery(ctx)
	_ = storage.ContextStorage(ctx, "any")
	_ = queue.ContextQueueService(ctx)
	_ = task.ContextTaskService(ctx)

	interceptor := UnaryServerInterceptor()
	var saw context.Context
	_, err := interceptor(base, nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(c context.Context, req interface{}) (interface{}, error) {
		saw = c
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if saw == nil || saw == base {
		t.Fatal("expected implementor context from unary interceptor")
	}
}

func TestStreamImplementorContextStable(t *testing.T) {
	base := context.Background()
	var wrapper *implementorStream
	err := StreamServerInterceptor()(nil, &fakeStream{ctx: base}, &grpc.StreamServerInfo{}, func(srv interface{}, ss grpc.ServerStream) error {
		wrapper = ss.(*implementorStream)
		c1 := ss.Context()
		c2 := ss.Context()
		if c1 != c2 {
			t.Fatal("Context() should return stable implementor context")
		}
		_ = client.ContextConnPool(c1)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if wrapper == nil {
		t.Fatal("expected implementor stream wrapper")
	}
}

type fakeStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (f *fakeStream) Context() context.Context { return f.ctx }
