package auth

import (
	"github.com/appootb/protobuf/go/service"
	"google.golang.org/grpc"
)

var (
	impl service.Authenticator
)

// Return the service implementor.
func Implementor() service.Authenticator {
	return impl
}

// Register service implementor.
func RegisterImplementor(auth service.Authenticator) {
	impl = auth
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return service.UnaryServerInterceptor(impl)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return service.StreamServerInterceptor(impl)
}
