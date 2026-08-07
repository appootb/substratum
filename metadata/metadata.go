// Package metadata parses client device and request metadata from gRPC incoming headers.
package metadata

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Well-known metadata header keys carried on RPC requests.
const (
	KeyProduct     = "prd"
	KeyPackage     = "pkg"
	KeyVersion     = "ver"
	KeyOSVersion   = "os"
	KeyBrand       = "brand"
	KeyModel       = "model"
	KeyDeviceID    = "udid"
	KeyFingerprint = "fp"
	KeyLocale      = "loc"
	KeyLatitude    = "lat"
	KeyLongitude   = "lon"
	KeyPlatform    = "plf"
	KeyNetwork     = "net"
	KeyTimestamp   = "time"
	KeyTraceID     = "sn"
	KeyIsEmulator  = "emulator"
	KeyIsDevelop   = "dev"
	KeyIsTesting   = "test"
	KeyChannel     = "channel"
	KeyUUID        = "uuid"
	KeyIMEI        = "imei"
	KeyDeviceMac   = "mac"
	KeyUserAgent   = "ua"
	KeyToken       = "token"

	KeyIANAUserAgent = "user-agent"
	KeyOriginalIP    = "x-forwarded-for"
	KeyAuthorization = "authorization"
)

var (
	// EnvDevelop gates whether the "dev" metadata flag is honored (when non-empty).
	EnvDevelop = os.Getenv("DEVELOPMENT")
	// EnvTesting gates whether the "test" metadata flag is honored (when non-empty).
	EnvTesting = os.Getenv("TESTING")
)

// mdFirst returns the first value for key, or "".
// gRPC lowercases keys; callers should pass lowercase constants.
func mdFirst(md metadata.MD, key string) string {
	vs := md[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// ParseIncomingMetadata builds common.Metadata from gRPC incoming metadata.
// Returns nil when the context has no incoming metadata.
//
// Prefer attaching via UnaryServerInterceptor / StreamServerInterceptor and
// reading with IncomingMetadata, which parses lazily once per request.
func ParseIncomingMetadata(ctx context.Context) *common.Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return parseFromMD(md, peerAddr(ctx))
}

// parseFromMD maps gRPC metadata values into common.Metadata.
// peerIP is used when x-forwarded-for is absent.
func parseFromMD(md metadata.MD, peerIP string) *common.Metadata {
	// Platform
	platform := common.Platform_PLATFORM_UNSPECIFIED
	if plf := mdFirst(md, KeyPlatform); plf != "" {
		if i, err := strconv.Atoi(plf); err != nil {
			platform = common.Platform(common.Platform_value[strings.ToUpper(plf)])
		} else {
			platform = common.Platform(i)
		}
	}
	// Network
	network := common.Network_NETWORK_UNSPECIFIED
	if net := mdFirst(md, KeyNetwork); net != "" {
		if i, err := strconv.Atoi(net); err != nil {
			network = common.Network(common.Network_value[strings.ToUpper(net)])
		} else {
			network = common.Network(i)
		}
	}
	// Timestamp
	timestamp := int64(0)
	if ts := mdFirst(md, KeyTimestamp); ts != "" {
		timestamp, _ = strconv.ParseInt(ts, 10, 64)
	}
	// Emulator
	emulator := false
	if v := mdFirst(md, KeyIsEmulator); v != "" {
		emulator, _ = strconv.ParseBool(v)
	}
	// Develop
	develop := false
	if EnvDevelop != "" {
		if v := mdFirst(md, KeyIsDevelop); v != "" {
			develop, _ = strconv.ParseBool(v)
		}
	}
	// Testing
	testing := false
	if EnvTesting != "" {
		if v := mdFirst(md, KeyIsTesting); v != "" {
			testing, _ = strconv.ParseBool(v)
		}
	}
	// Client IP
	clientIP := mdFirst(md, KeyOriginalIP)
	if clientIP == "" {
		clientIP = peerIP
	}
	// User agent
	userAgent := mdFirst(md, KeyUserAgent)
	if userAgent == "" {
		userAgent = mdFirst(md, KeyIANAUserAgent)
	}

	token := mdFirst(md, KeyToken)
	if token == "" {
		if auth := mdFirst(md, KeyAuthorization); strings.HasPrefix(auth, "Bearer ") {
			token = strings.TrimPrefix(auth, "Bearer ")
		}
	}

	return &common.Metadata{
		Product:     mdFirst(md, KeyProduct),
		Package:     mdFirst(md, KeyPackage),
		Version:     mdFirst(md, KeyVersion),
		OsVersion:   mdFirst(md, KeyOSVersion),
		Brand:       mdFirst(md, KeyBrand),
		Model:       mdFirst(md, KeyModel),
		DeviceId:    mdFirst(md, KeyDeviceID),
		Fingerprint: mdFirst(md, KeyFingerprint),
		Locale:      mdFirst(md, KeyLocale),
		Latitude:    mdFirst(md, KeyLatitude),
		Longitude:   mdFirst(md, KeyLongitude),
		Platform:    platform,
		Network:     network,
		Timestamp:   timestamp,
		TraceId:     mdFirst(md, KeyTraceID),
		IsEmulator:  emulator,
		IsDevelop:   develop,
		IsTesting:   testing,
		Channel:     mdFirst(md, KeyChannel),
		Uuid:        mdFirst(md, KeyUUID),
		Imei:        mdFirst(md, KeyIMEI),
		DeviceMac:   mdFirst(md, KeyDeviceMac),
		ClientIp:    clientIP,
		UserAgent:   userAgent,
		Token:       token,
	}
}

func peerAddr(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		return p.Addr.String()
	}
	return ""
}
