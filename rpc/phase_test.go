package rpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestBuildUnaryChainPhaseOrder(t *testing.T) {
	resetPhases()
	t.Cleanup(resetPhases)

	var order []string
	mark := func(name string) grpc.UnaryServerInterceptor {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			order = append(order, name)
			return handler(ctx, req)
		}
	}

	RegisterUnary(UnaryOutermost, mark("outermost"))
	RegisterUnary(UnaryBeforeAuth, mark("before_auth"))
	RegisterUnary(UnaryAfterAuth, mark("after_auth"))
	RegisterUnary(UnaryBeforeHandler, mark("before_handler"))

	chain := BuildUnaryChain()
	if len(chain) != 16 {
		t.Fatalf("chain length=%d want 16", len(chain))
	}

	// 4 framework interceptors sit between outermost and before_auth.
	if chain[0] == nil || chain[5] == nil || chain[7] == nil || chain[15] == nil {
		t.Fatal("custom phase interceptors missing at expected indices")
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		order = append(order, "handler")
		return nil, nil
	}
	// Nest only custom phase interceptors to verify relative order.
	chained := chain[0]
	for _, next := range []grpc.UnaryServerInterceptor{chain[5], chain[7], chain[15]} {
		cur := next
		prev := chained
		chained = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
			return prev(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
				return cur(ctx, req, info, h)
			})
		}
	}
	if _, err := chained(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler); err != nil {
		t.Fatal(err)
	}
	want := []string{"outermost", "before_auth", "after_auth", "before_handler", "handler"}
	if len(order) != len(want) {
		t.Fatalf("order=%v want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order=%v want %v", order, want)
		}
	}
}

func TestRegisterUnaryAfterFreezePanics(t *testing.T) {
	resetPhases()
	t.Cleanup(resetPhases)

	BuildUnaryChain()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when registering after freeze")
		}
	}()
	RegisterUnary(UnaryAfterAuth, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	})
}

func TestBuildUnaryChainIsIdempotent(t *testing.T) {
	resetPhases()
	t.Cleanup(resetPhases)

	first := BuildUnaryChain()
	second := BuildUnaryChain()
	if len(first) != len(second) {
		t.Fatalf("length mismatch: %d vs %d", len(first), len(second))
	}
	if &first[0] != &second[0] {
		t.Fatal("expected cached chain slice on second build")
	}
}

func TestBuildInterceptorChainsBuildsBoth(t *testing.T) {
	resetPhases()
	t.Cleanup(resetPhases)

	BuildInterceptorChains()
	if frozenUnaryChain == nil || frozenStreamChain == nil {
		t.Fatal("expected both chains built")
	}
	if len(frozenUnaryChain) != 12 {
		t.Fatalf("unary chain length=%d want 12", len(frozenUnaryChain))
	}
	if len(frozenStreamChain) != 12 {
		t.Fatalf("stream chain length=%d want 12", len(frozenStreamChain))
	}
}

func TestBuildStreamChainIncludesPhases(t *testing.T) {
	resetPhases()
	t.Cleanup(resetPhases)

	RegisterStream(StreamOutermost, func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	})
	RegisterStream(StreamAfterAuth, func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	})

	chain := BuildStreamChain()
	if len(chain) != 14 {
		t.Fatalf("stream chain length=%d want 14", len(chain))
	}
}
