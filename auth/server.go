package auth

import (
	"github.com/appootb/substratum/service"
	"google.golang.org/grpc"
)

var (
	impl service.Authenticator
)

// Implementor returns the authenticator service implementor.
func Implementor() service.Authenticator {
	return impl
}

// RegisterImplementor register the authenticator service implementor.
func RegisterImplementor(auth service.Authenticator) {
	impl = auth
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return service.UnaryServerInterceptor(impl)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return service.StreamServerInterceptor(impl)
}
