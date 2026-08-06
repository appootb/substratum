// Package recovery provides gRPC server interceptors that recover from panics
// and convert them into Internal status errors with a stack dump.
package recovery

import (
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a unary interceptor that recovers from panics.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) error {
		debug.PrintStack()
		return status.Errorf(codes.Internal, "%v", p)
	}))
}

// StreamServerInterceptor returns a stream interceptor that recovers from panics.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(func(p any) error {
		debug.PrintStack()
		return status.Errorf(codes.Internal, "%v", p)
	}))
}
