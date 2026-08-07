package logger

import (
	"context"
	"reflect"
	"testing"

	md "github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type captureLogger struct {
	lastMsg string
	last    Content
}

func (c *captureLogger) UpdateLevel(Level) {}
func (c *captureLogger) Log(_ Level, _ *common.Metadata, msg string, content Content) {
	c.lastMsg = msg
	c.last = content
}

func withCapture(t *testing.T) *captureLogger {
	t.Helper()
	c := &captureLogger{}
	prev := impl
	impl = c
	t.Cleanup(func() { impl = prev })
	return c
}

func withParsedMD(t *testing.T, pairs ...string) context.Context {
	t.Helper()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(pairs...))
	out, err := md.UnaryServerInterceptor()(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return ctx, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out.(context.Context)
}

func TestMetadataForLogDoesNotMutateSource(t *testing.T) {
	ctx := withParsedMD(t, md.KeyToken, "secret", md.KeyProduct, "app")
	src := md.IncomingMetadata(ctx)
	if src == nil || src.Token != "secret" {
		t.Fatalf("setup token=%v", src)
	}

	got := metadataForLog(ctx)
	if got == nil || got.Token != "" {
		t.Fatalf("log md token=%q want empty", got.GetToken())
	}
	if src.Token != "secret" {
		t.Fatalf("source token mutated: %q", src.Token)
	}
	if got.Product != "app" {
		t.Fatalf("product=%q", got.Product)
	}
}

func TestUnaryAccessLogIncludesPayloads(t *testing.T) {
	c := withCapture(t)
	SetAccessLogEnabled(true)
	t.Cleanup(func() { SetAccessLogEnabled(true) })

	ctx := withParsedMD(t, md.KeyProduct, "p")
	req := map[string]string{"k": "v"}
	resp := map[string]string{"ok": "1"}
	interceptor := UnaryServerInterceptor()
	_, err := interceptor(ctx, req, &grpc.UnaryServerInfo{FullMethod: "/svc/Method"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return resp, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.lastMsg != AccessLog {
		t.Fatalf("msg=%q", c.lastMsg)
	}
	if c.last[LogRequest] == nil || c.last[LogResponse] == nil {
		t.Fatalf("payloads=%v", c.last)
	}
	if !reflect.DeepEqual(c.last[LogRequest], req) || !reflect.DeepEqual(c.last[LogResponse], resp) {
		t.Fatalf("payloads=%v", c.last)
	}
	if c.last[LogPath] != "/svc/Method" {
		t.Fatalf("path=%v", c.last[LogPath])
	}
}

func TestSetAccessLogEnabledSkipsInterceptor(t *testing.T) {
	c := withCapture(t)
	SetAccessLogEnabled(false)
	t.Cleanup(func() { SetAccessLogEnabled(true) })

	if AccessLogEnabled() {
		t.Fatal("expected disabled")
	}
	ctx := withParsedMD(t, md.KeyProduct, "p")
	interceptor := UnaryServerInterceptor()
	_, err := interceptor(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/svc/Method"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.lastMsg != "" {
		t.Fatalf("expected no access log, got %q", c.lastMsg)
	}
}
