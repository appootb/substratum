package errors

import (
	"context"
	"net/url"
	"reflect"
	"strconv"
	"time"

	md "github.com/appootb/substratum/v2/metadata"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	impl Prompter
)

// Implementor returns the error prompter service implementor.
func Implementor() Prompter {
	return impl
}

// RegisterImplementor registers the error prompter service implementor.
func RegisterImplementor(prompter Prompter) {
	if impl != nil {
		panic("errors: implementor already registered")
	}
	impl = prompter
}

type Prompter interface {
	Translate(lang string, code int32) string
}

// UnaryResponseInterceptor rewrites errors for clients and attaches response
// metadata used by gRPC-Gateway (Appootb-code / message / timestamp).
//
// Success path sets timestamp only (clients use it for clock sync). code/message
// are omitted on success; gateway JSON envelopes already carry code=0 in the body.
// Error path sets code, message, and timestamp for HTTP/gRPC clients.
func UnaryResponseInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			// Return (typed-nil pointer, nil) can break protobuf marshal.
			if isNilPointerResponse(resp) {
				return reflect.New(reflect.TypeOf(resp).Elem()).Interface(), nil
			}
			_ = grpc.SetHeader(ctx, metadata.MD{
				"timestamp": {strconv.FormatInt(time.Now().Unix(), 10)},
			})
			return resp, nil
		}

		code, message := strconv.Itoa(int(codes.Unknown)), url.QueryEscape(err.Error())
		locale := "en"
		if incomingMD := md.IncomingMetadata(ctx); incomingMD != nil {
			if loc := incomingMD.GetLocale(); loc != "" {
				locale = loc
			}
		}

		if se, ok := err.(*StatusError); ok {
			if translated := Implementor().Translate(locale, se.Code); translated != "" {
				se.Message = translated
			}
			code = strconv.Itoa(int(se.Code))
			message = url.QueryEscape(se.Message)
			err = status.ErrorProto((*spb.Status)(se))
		} else if s, ok := status.FromError(err); ok {
			sp := s.Proto()
			if translated := Implementor().Translate(locale, sp.GetCode()); translated != "" {
				sp.Message = translated
			}
			code = strconv.Itoa(int(sp.GetCode()))
			message = url.QueryEscape(sp.GetMessage())
			err = status.ErrorProto(sp)
		}

		// Gateway/HTTP clients read these via OutgoingHeaderMatcher (Appootb-*).
		_ = grpc.SetHeader(ctx, metadata.MD{
			"code":      {code},
			"message":   {message},
			"timestamp": {strconv.FormatInt(time.Now().Unix(), 10)},
		})
		return resp, err
	}
}

// isNilPointerResponse reports whether resp is a typed nil pointer that would
// break protobuf marshaling. Untyped nil and non-pointer nilables are left alone.
func isNilPointerResponse(resp interface{}) bool {
	if resp == nil {
		return false
	}
	v := reflect.ValueOf(resp)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	// TODO
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
