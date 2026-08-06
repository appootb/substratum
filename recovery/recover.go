package recovery

import (
	"os"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const publicInternalMessage = "internal server error"

// RecoverStatus logs the panic stack and returns a client-safe gRPC status.
func RecoverStatus(p interface{}) error {
	debug.PrintStack()
	if os.Getenv("DEVELOPMENT") != "" {
		return status.Errorf(codes.Internal, "%v", p)
	}
	return status.Error(codes.Internal, publicInternalMessage)
}
