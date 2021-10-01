package recovery

import (
	"runtime/debug"

	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p interface{}) error {
		debug.PrintStack()
		return status.Errorf(codes.Internal, "%v", p)
	}))
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(func(p interface{}) error {
		debug.PrintStack()
		return status.Errorf(codes.Internal, "%v", p)
	}))
}
