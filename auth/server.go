package auth

import (
	"github.com/appootb/protobuf/go/service"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return service.UnaryServerInterceptor(Default)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return service.StreamServerInterceptor(Default)
}
