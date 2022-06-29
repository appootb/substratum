package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	md "github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/util/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	MetadataHeaderPrefix  = "Appootb-"
	MetadataTrailerPrefix = "Trailer-"
)

var DefaultOptions = []runtime.ServeMuxOption{
	runtime.WithMarshalerOption(MIMEJSON, &JSONMarshal{}),
	runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
	runtime.WithOutgoingHeaderMatcher(OutgoingHeaderMatcher),
	runtime.WithMetadata(URLQueryMetadata),
	runtime.WithErrorHandler(ProtoErrorHandler),
	runtime.WithStreamErrorHandler(StreamErrorHandler),
}

func New(opts []runtime.ServeMuxOption) *runtime.ServeMux {
	return runtime.NewServeMux(opts...)
}

// isPermanentHTTPHeader checks whether hdr belongs to the list of
// permanent request headers maintained by IANA.
// http://www.iana.org/assignments/message-headers/message-headers.xml
func isPermanentHTTPHeader(hdr string) bool {
	switch hdr {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Language",
		"Accept-Ranges",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Cookie",
		"Date",
		"Expect",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Schedule-Tag-Match",
		"If-Unmodified-Since",
		"Max-Forwards",
		"Origin",
		"Pragma",
		"Referer",
		"User-Agent",
		"Via",
		"Warning":
		return true
	}
	return false
}

func IncomingHeaderMatcher(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	if isPermanentHTTPHeader(key) {
		return key, true
	} else if strings.HasPrefix(key, MetadataHeaderPrefix) {
		return key[len(MetadataHeaderPrefix):], true
	}
	return "", false
}

func OutgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", MetadataHeaderPrefix, key), true
}

func URLQueryMetadata(_ context.Context, r *http.Request) metadata.MD {
	queryMD := make(metadata.MD, len(r.URL.Query()))
	for k, v := range r.URL.Query() {
		switch k {
		case md.KeyProduct, md.KeyPackage, md.KeyVersion,
			md.KeyOSVersion, md.KeyBrand, md.KeyModel, md.KeyDeviceID, md.KeyFingerprint,
			md.KeyLocale, md.KeyLatitude, md.KeyLongitude,
			md.KeyPlatform, md.KeyNetwork,
			md.KeyTimestamp, md.KeyTraceID,
			md.KeyIsEmulator, md.KeyIsDevelop, md.KeyIsTesting,
			md.KeyChannel, md.KeyUUID, md.KeyIMEI, md.KeyDeviceMac, md.KeyUserAgent,
			md.KeyToken:
			queryMD[k] = v
		}
	}
	return queryMD
}

func ProtoErrorHandler(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error) {
	//
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType(nil))

	ctxMD, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		// TODO
	}
	// Metadata
	for k, vs := range ctxMD.HeaderMD {
		nk := fmt.Sprintf("%s%s", MetadataHeaderPrefix, k)
		for _, v := range vs {
			w.Header().Add(nk, v)
		}
	}
	// Trailer header
	for k := range ctxMD.TrailerMD {
		tk := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tk)
	}
	// Write response
	body := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
	}
	if s, ok := status.FromError(err); ok {
		body = s.Proto()
	}
	buf, err := jsonpb.Marshal(body)
	if err != nil {
		buf = []byte(`{"error": "failed to marshal error message"}`)
	}
	if _, err = w.Write(buf); err != nil {
		// TODO
	}

	// Trailer
	for k, vs := range ctxMD.TrailerMD {
		tk := fmt.Sprintf("%s%s", MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tk, v)
		}
	}
}

func StreamErrorHandler(_ context.Context, err error) *status.Status {
	// TODO
	return status.Convert(err)
}
