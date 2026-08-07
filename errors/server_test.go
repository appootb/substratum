package errors

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type emptyPrompter struct{}

func (emptyPrompter) Translate(string, int32) string { return "" }

func ensurePrompter(t *testing.T) {
	t.Helper()
	if Implementor() == nil {
		RegisterImplementor(emptyPrompter{})
	}
}

func TestUnarySuccessReturnsResponse(t *testing.T) {
	ensurePrompter(t)
	interceptor := UnaryResponseInterceptor()
	want := &emptypb.Empty{}
	got, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return want, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got=%v want same pointer", got)
	}
}

func TestUnarySuccessRepairsTypedNilPointer(t *testing.T) {
	ensurePrompter(t)
	interceptor := UnaryResponseInterceptor()
	var nilEmpty *emptypb.Empty
	got, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nilEmpty, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected non-nil empty message")
	}
	if _, ok := got.(*emptypb.Empty); !ok {
		t.Fatalf("type=%T", got)
	}
}

func TestUnarySuccessDoesNotRepairNilMap(t *testing.T) {
	ensurePrompter(t)
	interceptor := UnaryResponseInterceptor()
	var nilMap map[string]string
	got, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nilMap, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	m, ok := got.(map[string]string)
	if !ok {
		t.Fatalf("type=%T", got)
	}
	if m != nil {
		t.Fatalf("expected nil map passthrough, got %v", m)
	}
}

func TestUnaryErrorTranslatesStatus(t *testing.T) {
	ensurePrompter(t)
	interceptor := UnaryResponseInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, New(codes.NotFound, "missing")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Fatalf("status=%v", err)
	}
}

func TestIsNilPointerResponse(t *testing.T) {
	var p *emptypb.Empty
	if !isNilPointerResponse(p) {
		t.Fatal("typed nil pointer should match")
	}
	if isNilPointerResponse(nil) {
		t.Fatal("untyped nil should not match")
	}
	var m map[string]string
	if isNilPointerResponse(m) {
		t.Fatal("nil map should not match")
	}
}
