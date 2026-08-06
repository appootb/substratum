// Package monitor exposes gRPC server Prometheus metrics interceptors.
package monitor

import (
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
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return serverMetrics.UnaryServerInterceptor()
}

// StreamServerInterceptor returns the Prometheus stream interceptor.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return serverMetrics.StreamServerInterceptor()
}

// InitializeMetrics pre-creates metric series for all methods registered on server.
// Call after service registration (equivalent to go-grpc-prometheus.Register).
func InitializeMetrics(server *grpc.Server) {
	if server == nil {
		return
	}
	serverMetrics.InitializeMetrics(server)
}
