package metadata

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestLazyIncomingMetadataParsesOnce(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		KeyProduct, "app",
		KeyToken, "tok",
	))
	ctx = ContextWithIncomingMetadata(ctx)

	h, ok := ctx.Value(metadataKey{}).(*holder)
	if !ok || h == nil {
		t.Fatal("expected holder")
	}
	if h.parsed != nil {
		t.Fatal("should not parse before IncomingMetadata")
	}

	a := IncomingMetadata(ctx)
	b := IncomingMetadata(ctx)
	if a == nil || a.Product != "app" || a.Token != "tok" {
		t.Fatalf("parsed=%+v", a)
	}
	if a != b {
		t.Fatal("expected same cached pointer")
	}
	if h.parsed != a {
		t.Fatal("holder cache mismatch")
	}
}

func TestUnaryInterceptorDefersParse(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(KeyProduct, "p"))
	var saw context.Context
	_, err := UnaryServerInterceptor()(ctx, nil, &grpc.UnaryServerInfo{}, func(c context.Context, req interface{}) (interface{}, error) {
		saw = c
		h := c.Value(metadataKey{}).(*holder)
		if h.parsed != nil {
			t.Fatal("parsed before read")
		}
		if IncomingMetadata(c).Product != "p" {
			t.Fatal("product")
		}
		if h.parsed == nil {
			t.Fatal("expected cache after read")
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if saw == nil {
		t.Fatal("handler not called")
	}
}

func TestStreamWrapperReusesContext(t *testing.T) {
	base := metadata.NewIncomingContext(context.Background(), metadata.Pairs(KeyLocale, "en"))
	var wrapper *ctxWrapper
	err := StreamServerInterceptor()(nil, &fakeStream{ctx: base}, &grpc.StreamServerInfo{}, func(srv interface{}, ss grpc.ServerStream) error {
		wrapper = ss.(*ctxWrapper)
		c1 := ss.Context()
		c2 := ss.Context()
		if c1 != c2 {
			t.Fatal("Context() should return stable ctx")
		}
		a := IncomingMetadata(c1)
		b := IncomingMetadata(c2)
		if a != b || a.Locale != "en" {
			t.Fatalf("md=%v %v", a, b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if wrapper == nil {
		t.Fatal("no wrapper")
	}
}

func TestContextWithIncomingMetadataIdempotent(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(KeyProduct, "x"))
	ctx1 := ContextWithIncomingMetadata(ctx)
	ctx2 := ContextWithIncomingMetadata(ctx1)
	if ctx1 != ctx2 {
		t.Fatal("expected same context when already attached")
	}
}

func TestIncomingMetadataMissing(t *testing.T) {
	if IncomingMetadata(context.Background()) != nil {
		t.Fatal("expected nil")
	}
	// No incoming MD: attach is a no-op.
	ctx := ContextWithIncomingMetadata(context.Background())
	if IncomingMetadata(ctx) != nil {
		t.Fatal("expected nil without incoming metadata")
	}
}

type fakeStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (f *fakeStream) Context() context.Context { return f.ctx }
