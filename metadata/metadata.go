package metadata

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

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
)

var (
	EnvDevelop = os.Getenv("DEVELOPMENT")
	EnvTesting = os.Getenv("TESTING")
)

func ParseIncomingMetadata(ctx context.Context) *common.Metadata {
	var md url.Values
	if v, ok := metadata.FromIncomingContext(ctx); !ok {
		return nil
	} else {
		md = url.Values(v)
	}
	// Platform
	platform := common.Platform_PLATFORM_UNSPECIFIED
	if md.Has(KeyPlatform) {
		if i, err := strconv.Atoi(md.Get(KeyPlatform)); err != nil {
			platform = common.Platform(common.Platform_value[strings.ToUpper(md.Get(KeyPlatform))])
		} else {
			platform = common.Platform(i)
		}
	}
	// Network
	network := common.Network_NETWORK_UNSPECIFIED
	if md.Has(KeyNetwork) {
		if i, err := strconv.Atoi(md.Get(KeyNetwork)); err != nil {
			network = common.Network(common.Network_value[strings.ToUpper(md.Get(KeyNetwork))])
		} else {
			network = common.Network(i)
		}
	}
	// Timestamp
	timestamp := int64(0)
	if md.Has(KeyTimestamp) {
		timestamp, _ = strconv.ParseInt(md.Get(KeyTimestamp), 10, 64)
	}
	// Emulator
	emulator := false
	if md.Has(KeyIsEmulator) {
		emulator, _ = strconv.ParseBool(md.Get(KeyIsEmulator))
	}
	// Develop
	develop := false
	if md.Has(KeyIsDevelop) && EnvDevelop != "" {
		develop, _ = strconv.ParseBool(md.Get(KeyIsDevelop))
	}
	// Testing
	testing := false
	if md.Has(KeyIsTesting) && EnvTesting != "" {
		testing, _ = strconv.ParseBool(md.Get(KeyIsTesting))
	}
	// Client IP
	clientIP := ""
	if md.Has(KeyOriginalIP) {
		clientIP = md.Get(KeyOriginalIP)
	} else if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	// User agent
	userAgent := ""
	if md.Has(KeyUserAgent) {
		userAgent = md.Get(KeyUserAgent)
	} else {
		userAgent = md.Get(KeyIANAUserAgent)
	}

	return &common.Metadata{
		Product:     md.Get(KeyProduct),
		Package:     md.Get(KeyPackage),
		Version:     md.Get(KeyVersion),
		OsVersion:   md.Get(KeyOSVersion),
		Brand:       md.Get(KeyBrand),
		Model:       md.Get(KeyModel),
		DeviceId:    md.Get(KeyDeviceID),
		Fingerprint: md.Get(KeyFingerprint),
		Locale:      md.Get(KeyLocale),
		Latitude:    md.Get(KeyLatitude),
		Longitude:   md.Get(KeyLongitude),
		Platform:    platform,
		Network:     network,
		Timestamp:   timestamp,
		TraceId:     md.Get(KeyTraceID),
		IsEmulator:  emulator,
		IsDevelop:   develop,
		IsTesting:   testing,
		Channel:     md.Get(KeyChannel),
		Uuid:        md.Get(KeyUUID),
		Imei:        md.Get(KeyIMEI),
		DeviceMac:   md.Get(KeyDeviceMac),
		ClientIp:    clientIP,
		UserAgent:   userAgent,
		Token:       md.Get(KeyToken),
	}
}
