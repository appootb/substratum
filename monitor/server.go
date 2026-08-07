// Package monitor exposes gRPC server Prometheus metrics interceptors.
package monitor

import (
	"context"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// serverMetrics is the process-wide gRPC server metrics collector
// (replaces the deprecated go-grpc-prometheus DefaultServerMetrics singleton).
var serverMetrics = func() *grpcprom.ServerMetrics {
	m := grpcprom.NewServerMetrics()
	prometheus.MustRegister(m)
	return m
}()

// UnaryServerInterceptor returns the Prometheus unary interceptor.
// When metrics are disabled (see SetMetricsEnabled), the interceptor is a no-op.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	prom := serverMetrics.UnaryServerInterceptor()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !MetricsEnabled() {
			return handler(ctx, req)
		}
		return prom(ctx, req, info, handler)
	}
}

// StreamServerInterceptor returns the Prometheus stream interceptor.
// When metrics are disabled (see SetMetricsEnabled), the interceptor is a no-op.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	prom := serverMetrics.StreamServerInterceptor()
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !MetricsEnabled() {
			return handler(srv, ss)
		}
		return prom(srv, ss, info, handler)
	}
}

// InitializeMetrics pre-creates metric series for all methods registered on server.
// Call after service registration (equivalent to go-grpc-prometheus.Register).
// No-op when metrics are disabled.
func InitializeMetrics(server *grpc.Server) {
	if server == nil || !MetricsEnabled() {
		return
	}
	serverMetrics.InitializeMetrics(server)
}
