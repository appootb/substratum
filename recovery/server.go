package recovery

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = toPanicError(r, req, info.FullMethod)
			}
		}()
		return handler(ctx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = toPanicError(r, "", info.FullMethod)
			}
		}()
		return handler(srv, stream)
	}
}

func toPanicError(r, req interface{}, m string) error {
	fmt.Printf("path: %s\t request:%+v\tpanic:%+v\n%s", m, req, r, debug.Stack())
	return status.Errorf(codes.Internal, "%v", r)
}
