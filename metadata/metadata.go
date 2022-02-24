package metadata

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/appootb/substratum/proto/go/common"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	KeyToken      = "token"
	KeyPlatform   = "platform"
	KeyNetwork    = "network"
	KeyPackage    = "package"
	KeyVersion    = "version"
	KeyOSVersion  = "os_version"
	KeyBrand      = "brand"
	KeyModel      = "model"
	KeyDeviceID   = "device_id"
	KeyTimestamp  = "timestamp"
	KeyIsEmulator = "is_emulator"
	KeyIsDebug    = "is_debug"
	KeyLatitude   = "latitude"
	KeyLongitude  = "longitude"
	KeyLocale     = "locale"
	KeyChannel    = "channel"
	KeyProduct    = "product"
	KeyTraceID    = "trace_id"
	KeyRiskID     = "risk_id"
	KeyUUID       = "uuid"
	KeyUDID       = "udid"
	KeyUserAgent  = "user_agent"
	KeyDeviceMac  = "device_mac"
	KeyAndroidID  = "android_id"
	KeyOriginalIP = "x-forwarded-for"
)

var (
	Debug = os.Getenv("DEBUG")
)

func ParseIncomingMetadata(ctx context.Context) *common.Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	clientIP := ""
	if ips := md.Get(KeyOriginalIP); len(ips) > 0 {
		clientIP = ips[0]
	} else if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	emulator := false
	if b := md.Get(KeyIsEmulator); len(b) > 0 {
		emulator, _ = strconv.ParseBool(b[0])
	}
	platform := common.Platform_PLATFORM_UNSPECIFIED
	if pf := md.Get(KeyPlatform); len(pf) > 0 {
		if i, err := strconv.Atoi(pf[0]); err != nil {
			platform = common.Platform(common.Platform_value[pf[0]])
		} else {
			platform = common.Platform(i)
		}
	}
	timestamp := int64(0)
	if ts := md.Get(KeyTimestamp); len(ts) > 0 {
		timestamp, _ = strconv.ParseInt(ts[0], 10, 64)
	}
	network := common.Network_NETWORK_UNSPECIFIED
	if net := md.Get(KeyNetwork); len(net) > 0 {
		if i, err := strconv.Atoi(net[0]); err != nil {
			network = common.Network(common.Network_value[net[0]])
		} else {
			network = common.Network(i)
		}
	}
	debug := false
	if b := md.Get(KeyIsDebug); len(b) > 0 && Debug != "" {
		debug, _ = strconv.ParseBool(b[0])
	}

	return &common.Metadata{
		Token:      strings.Join(md[KeyToken], ""),
		Package:    strings.Join(md[KeyPackage], ""),
		Version:    strings.Join(md[KeyVersion], ""),
		OsVersion:  strings.Join(md[KeyOSVersion], ""),
		Brand:      strings.Join(md[KeyBrand], ""),
		Model:      strings.Join(md[KeyModel], ""),
		UserAgent:  strings.Join(md[KeyUserAgent], ""),
		DeviceId:   strings.Join(md[KeyDeviceID], ""),
		Platform:   platform,
		Timestamp:  timestamp,
		IsEmulator: emulator,
		IsDebug:    debug,
		Network:    network,
		ClientIp:   clientIP,
		DeviceMac:  strings.Join(md[KeyDeviceMac], ""),
		Latitude:   strings.Join(md[KeyLatitude], ""),
		Longitude:  strings.Join(md[KeyLongitude], ""),
		Locale:     strings.Join(md[KeyLocale], ""),
		Channel:    strings.Join(md[KeyChannel], ""),
		Product:    strings.Join(md[KeyProduct], ""),
		TraceId:    strings.Join(md[KeyTraceID], ""),
		RiskId:     strings.Join(md[KeyRiskID], ""),
		Uuid:       strings.Join(md[KeyUUID], ""),
		Udid:       strings.Join(md[KeyUDID], ""),
		AndroidId:  strings.Join(md[KeyAndroidID], ""),
	}
}
