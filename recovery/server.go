package recovery

import (
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: support more.

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p interface{}) error {
		// TODO log
		return status.Errorf(codes.Internal, "%v", p)
	}))
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	// TODO log
	return recovery.StreamServerInterceptor()
}
