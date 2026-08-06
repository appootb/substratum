package rpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestChainUnaryServerOrder(t *testing.T) {
	var order []string
	mk := func(name string) grpc.UnaryServerInterceptor {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			order = append(order, name)
			return handler(ctx, req)
		}
	}
	chained := chainUnaryServer(mk("a"), mk("b"), mk("c"))
	_, err := chained(context.Background(), nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		order = append(order, "h")
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a", "b", "c", "h"}
	if len(order) != len(want) {
		t.Fatalf("order=%v want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order=%v want %v", order, want)
		}
	}
}

func TestChainUnaryServerEmpty(t *testing.T) {
	if chainUnaryServer() != nil {
		t.Fatal("expected nil chain")
	}
}

func TestChainStreamServerOrder(t *testing.T) {
	var order []string
	mk := func(name string) grpc.StreamServerInterceptor {
		return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			order = append(order, name)
			return handler(srv, ss)
		}
	}
	chained := chainStreamServer(mk("a"), mk("b"))
	err := chained(nil, nil, &grpc.StreamServerInfo{}, func(srv interface{}, stream grpc.ServerStream) error {
		order = append(order, "h")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a", "b", "h"}
	if len(order) != len(want) {
		t.Fatalf("order=%v", order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order=%v want %v", order, want)
		}
	}
}
