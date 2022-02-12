package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"

	"github.com/appootb/substratum/util/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var DefaultOptions = []runtime.ServeMuxOption{
	runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb.Marshaler),
	runtime.WithMetadata(URLQueryMetadata),
	runtime.WithProtoErrorHandler(ProtoErrorHandler),
	runtime.WithStreamErrorHandler(StreamErrorHandler),
}

func New(opts []runtime.ServeMuxOption) *runtime.ServeMux {
	return runtime.NewServeMux(opts...)
}

func URLQueryMetadata(_ context.Context, r *http.Request) metadata.MD {
	return metadata.MD(r.URL.Query())
}

func ProtoErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		// TODO
	}
	// Metadata
	for k, vs := range md.HeaderMD {
		nk := fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, k)
		for _, v := range vs {
			w.Header().Add(nk, v)
		}
	}
	// Trailer header
	for k := range md.TrailerMD {
		tk := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tk)
	}
	// Write response
	w.WriteHeader(http.StatusInternalServerError)
	body := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
	}
	if s, ok := status.FromError(err); ok {
		body = s.Proto()
	}
	buf, err := marshaler.Marshal(body)
	if err != nil {
		buf = []byte(`{"error": "failed to marshal error message"}`)
	}
	if _, err := w.Write(buf); err != nil {
		// TODO
	}

	// Trailer
	for k, vs := range md.TrailerMD {
		tk := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tk, v)
		}
	}
}

func StreamErrorHandler(ctx context.Context, err error) *runtime.StreamError {
	// TODO
	code := codes.Unknown
	message := err.Error()
	var details []*anypb.Any
	if s, ok := status.FromError(err); ok {
		code = s.Code()
		message = s.Message()
		details = s.Proto().GetDetails()
	}

	return &runtime.StreamError{
		GrpcCode:   int32(code),
		HttpCode:   http.StatusInternalServerError,
		Message:    message,
		HttpStatus: http.StatusText(http.StatusInternalServerError),
		Details:    details,
	}
}
