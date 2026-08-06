package recovery

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptorRecovers(t *testing.T) {
	interceptor := UnaryServerInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("boom")
		})
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Internal {
		t.Fatalf("status=%v ok=%v", err, ok)
	}
}

func TestUnaryServerInterceptorPassthrough(t *testing.T) {
	interceptor := UnaryServerInterceptor()
	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/t"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		})
	if err != nil || resp != "ok" {
		t.Fatalf("resp=%v err=%v", resp, err)
	}
}
