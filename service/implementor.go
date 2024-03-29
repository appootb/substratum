package service

import (
	"context"
	"net/http"

	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HttpHandler interface.
type HttpHandler interface {
	// Handle registers the handler for the given pattern.
	Handle(pattern string, handler http.Handler)

	// HandleFunc registers the handler function for the given pattern.
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Implementor interface.
type Implementor interface {
	// Context returns the service context.
	Context() context.Context

	// UnaryInterceptor returns the unary server interceptor for local gateway handler server.
	UnaryInterceptor() grpc.UnaryServerInterceptor

	// StreamInterceptor returns the stream server interceptor for local gateway handler server.
	StreamInterceptor() grpc.StreamServerInterceptor

	// GetGRPCServer returns gRPC server of the specified visible scope.
	GetGRPCServer(scope permission.VisibleScope) []*grpc.Server

	// GetGatewayMux returns gateway mux of the specified visible scope.
	GetGatewayMux(scope permission.VisibleScope) []*runtime.ServeMux
}

// Authenticator interface.
type Authenticator interface {
	// ServiceComponentName returns the component name implements the service method.
	ServiceComponentName(serviceMethod string) string

	// RegisterServiceSubjects registers required method subjects of the service.
	// The map key of the parameter is the full url path of the method.
	RegisterServiceSubjects(component string, serviceMethodSubjects map[string][]permission.Subject, serviceMethodRoles map[string][]string)

	// Authenticate a request specified by the full url path of the method.
	Authenticate(ctx context.Context, serviceMethod string) (*secret.Info, error)
}

type Method struct {
	// FullMethod is the full service method string.
	FullMethod string
	// IsClientStream indicates whether the service is a client streaming RPC.
	IsClientStream bool
	// IsServerStream indicates whether the service is a server streaming RPC.
	IsServerStream bool
	// IsHttpGateway indicates whether the service is a http request.
	IsHttpGateway bool
}

type serviceMethodKey struct{}

type componentKey struct{}

type secretKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor that authenticates incoming messages.
//
// Invalid messages will be rejected with `PermissionDenied` before reaching any userspace handlers.
func UnaryServerInterceptor(v Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		secretInfo, err := v.Authenticate(ctx, info.FullMethod)
		if err != nil {
			if _, ok := err.(interface {
				GRPCStatus() *status.Status
			}); ok {
				return nil, err
			}
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		return handler(context.WithValue(
			ContextWithComponentName(
				ContextWithServiceMethod(ctx, &Method{
					FullMethod: info.FullMethod,
				}),
				v.ServiceComponentName(info.FullMethod)),
			secretKey{}, secretInfo), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that authenticates incoming messages.
//
// The stage at which unauthenticated messages will be rejected with `PermissionDenied` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any user space
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamServerInterceptor(v Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		secretInfo, err := v.Authenticate(stream.Context(), info.FullMethod)
		if err != nil {
			if _, ok := err.(interface {
				GRPCStatus() *status.Status
			}); ok {
				return err
			}
			return status.Errorf(codes.PermissionDenied, err.Error())
		}
		wrapper := &ctxWrapper{
			ServerStream:     stream,
			StreamServerInfo: info,
			secret:           secretInfo,
			component:        v.ServiceComponentName(info.FullMethod),
		}
		return handler(srv, wrapper)
	}
}

type ctxWrapper struct {
	grpc.ServerStream
	*grpc.StreamServerInfo
	secret    *secret.Info
	component string
}

func (s *ctxWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()
	return context.WithValue(
		ContextWithComponentName(
			ContextWithServiceMethod(ctx, &Method{
				FullMethod:     s.StreamServerInfo.FullMethod,
				IsClientStream: s.StreamServerInfo.IsClientStream,
				IsServerStream: s.StreamServerInfo.IsServerStream,
			}),
			s.component),
		secretKey{}, s.secret)
}

func ContextWithServiceMethod(ctx context.Context, method *Method) context.Context {
	return context.WithValue(ctx, serviceMethodKey{}, method)
}

func ServiceMethodFromContext(ctx context.Context) *Method {
	if method := ctx.Value(serviceMethodKey{}); method != nil {
		return method.(*Method)
	}
	return &Method{}
}

func ContextWithComponentName(ctx context.Context, component string) context.Context {
	return context.WithValue(ctx, componentKey{}, component)
}

func ComponentNameFromContext(ctx context.Context) string {
	if component := ctx.Value(componentKey{}); component != nil {
		return component.(string)
	}
	return ""
}

func AccountSecretFromContext(ctx context.Context) *secret.Info {
	if secretInfo := ctx.Value(secretKey{}); secretInfo != nil {
		return secretInfo.(*secret.Info)
	}
	return &secret.Info{}
}
