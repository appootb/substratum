// Package recovery provides gRPC server interceptors that recover from panics
// and convert them into Internal status errors with a stack dump.
package recovery

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a unary interceptor that recovers from panics.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(RecoverStatus))
}

// StreamServerInterceptor returns a stream interceptor that recovers from panics.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(RecoverStatus))
}
