package monitor

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestSetMetricsEnabledDefaultTrue(t *testing.T) {
	// Restore default after any prior test mutation in this package.
	SetMetricsEnabled(true)
	t.Cleanup(func() { SetMetricsEnabled(true) })
	if !MetricsEnabled() {
		t.Fatal("expected metrics enabled by default")
	}
}

func TestSetMetricsEnabledPassthrough(t *testing.T) {
	SetMetricsEnabled(false)
	t.Cleanup(func() { SetMetricsEnabled(true) })

	called := false
	interceptor := UnaryServerInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/t"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		called = true
		return "ok", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("handler not called")
	}
	if MetricsEnabled() {
		t.Fatal("expected disabled")
	}
}
