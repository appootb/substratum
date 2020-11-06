package errors

import (
	"context"
	"net/url"
	"reflect"
	"strconv"
	"time"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	impl Prompter
)

// Return the service implementor.
func Implementor() Prompter {
	return impl
}

// Register service implementor.
func RegisterImplementor(prompter Prompter) {
	impl = prompter
}

type Prompter interface {
	Translate(code int32) string
}

func UnaryResponseInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		outgoingMD := metadata.MD{
			"code":    []string{"0"},
			"message": []string{""},
		}
		defer func() {
			outgoingMD.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
			_ = grpc.SetHeader(ctx, outgoingMD)
		}()

		resp, err := handler(ctx, req)
		if err == nil {
			// Return (nil, nil) might cause pb marshal fail.
			if reflect.ValueOf(resp).IsNil() {
				return reflect.New(reflect.TypeOf(resp).Elem()).Interface(), nil
			}
			return resp, err
		}

		// Update error message.
		if se, ok := err.(*StatusError); ok {
			if message := Implementor().Translate(se.Code); message != "" {
				se.Message = message
			}
			outgoingMD.Set("code", strconv.Itoa(int(se.Code)))
			outgoingMD.Set("message", url.QueryEscape(se.Message))
			return resp, status.ErrorProto((*spb.Status)(se))
		} else if s, ok := status.FromError(err); ok {
			message := Implementor().Translate(int32(s.Code()))
			if message == "" {
				message = s.Message()
			}
			outgoingMD.Set("code", strconv.Itoa(int(s.Code())))
			outgoingMD.Set("message", url.QueryEscape(message))
		}
		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	// TODO
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
