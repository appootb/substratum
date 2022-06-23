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

const (
	KeyProduct     = "product"
	KeyPackage     = "package"
	KeyVersion     = "version"
	KeyOSVersion   = "os"
	KeyBrand       = "brand"
	KeyModel       = "model"
	KeyDeviceID    = "udid"
	KeyFingerprint = "fingerprint"
	KeyLocale      = "locale"
	KeyLatitude    = "lat"
	KeyLongitude   = "lon"
	KeyPlatform    = "platform"
	KeyNetwork     = "network"
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

	KeyOriginalIP = "x-forwarded-for"
)

var (
	EnvDevelop = os.Getenv("DEVELOPMENT")
	EnvTesting = os.Getenv("TESTING")
)

func ParseIncomingMetadata(ctx context.Context) *common.Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	// Platform
	platform := common.Platform_PLATFORM_UNSPECIFIED
	if pf := md.Get(KeyPlatform); len(pf) > 0 {
		if i, err := strconv.Atoi(pf[0]); err != nil {
			platform = common.Platform(common.Platform_value[pf[0]])
		} else {
			platform = common.Platform(i)
		}
	}
	// Network
	network := common.Network_NETWORK_UNSPECIFIED
	if net := md.Get(KeyNetwork); len(net) > 0 {
		if i, err := strconv.Atoi(net[0]); err != nil {
			network = common.Network(common.Network_value[net[0]])
		} else {
			network = common.Network(i)
		}
	}
	// Timestamp
	timestamp := int64(0)
	if ts := md.Get(KeyTimestamp); len(ts) > 0 {
		timestamp, _ = strconv.ParseInt(ts[0], 10, 64)
	}
	// Emulator
	emulator := false
	if b := md.Get(KeyIsEmulator); len(b) > 0 {
		emulator, _ = strconv.ParseBool(b[0])
	}
	// Develop
	develop := false
	if b := md.Get(KeyIsDevelop); len(b) > 0 && EnvDevelop != "" {
		develop, _ = strconv.ParseBool(b[0])
	}
	// Testing
	testing := false
	if b := md.Get(KeyIsTesting); len(b) > 0 && EnvTesting != "" {
		testing, _ = strconv.ParseBool(b[0])
	}
	// Client IP
	clientIP := ""
	if ips := md.Get(KeyOriginalIP); len(ips) > 0 {
		clientIP = ips[0]
	} else if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}

	return &common.Metadata{
		Product:     strings.Join(md[KeyProduct], ""),
		Package:     strings.Join(md[KeyPackage], ""),
		Version:     strings.Join(md[KeyVersion], ""),
		OsVersion:   strings.Join(md[KeyOSVersion], ""),
		Brand:       strings.Join(md[KeyBrand], ""),
		Model:       strings.Join(md[KeyModel], ""),
		DeviceId:    strings.Join(md[KeyDeviceID], ""),
		Fingerprint: strings.Join(md[KeyFingerprint], ""),
		Locale:      strings.Join(md[KeyLocale], ""),
		Latitude:    strings.Join(md[KeyLatitude], ""),
		Longitude:   strings.Join(md[KeyLongitude], ""),
		Platform:    platform,
		Network:     network,
		Timestamp:   timestamp,
		TraceId:     strings.Join(md[KeyTraceID], ""),
		IsEmulator:  emulator,
		IsDevelop:   develop,
		IsTesting:   testing,
		Channel:     strings.Join(md[KeyChannel], ""),
		Uuid:        strings.Join(md[KeyUUID], ""),
		Imei:        strings.Join(md[KeyIMEI], ""),
		DeviceMac:   strings.Join(md[KeyDeviceMac], ""),
		ClientIp:    clientIP,
		UserAgent:   strings.Join(md[KeyUserAgent], ""),
		Token:       strings.Join(md[KeyToken], ""),
	}
}
