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

var (
	Debug = os.Getenv("DEBUG")
)

func ParseIncomingMetadata(ctx context.Context) *common.Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	clientIP := ""
	if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
		clientIP = ips[0]
	} else if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	emulator := true
	if b := md.Get("is_emulator"); len(b) > 0 {
		emulator, _ = strconv.ParseBool(b[0])
	}
	platform := common.Platform_PLATFORM_UNSPECIFIED
	if pf := md.Get("platform"); len(pf) > 0 {
		if i, err := strconv.Atoi(pf[0]); err != nil {
			platform = common.Platform(common.Platform_value[pf[0]])
		} else {
			platform = common.Platform(i)
		}
	}
	timestamp := int64(0)
	if ts := md.Get("timestamp"); len(ts) > 0 {
		timestamp, _ = strconv.ParseInt(ts[0], 10, 64)
	}
	network := common.Network_NETWORK_UNSPECIFIED
	if net := md.Get("network"); len(net) > 0 {
		if i, err := strconv.Atoi(net[0]); err != nil {
			network = common.Network(common.Network_value[net[0]])
		} else {
			network = common.Network(i)
		}
	}
	debug := false
	if b := md.Get("is_debug"); len(b) > 0 && Debug != "" {
		debug, _ = strconv.ParseBool(b[0])
	}

	return &common.Metadata{
		Token:      strings.Join(md["token"], ""),
		Package:    strings.Join(md["package"], ""),
		Version:    strings.Join(md["version"], ""),
		OsVersion:  strings.Join(md["os_version"], ""),
		Brand:      strings.Join(md["brand"], ""),
		Model:      strings.Join(md["model"], ""),
		UserAgent:  strings.Join(md["user_agent"], ""),
		DeviceId:   strings.Join(md["device_id"], ""),
		Platform:   platform,
		Timestamp:  timestamp,
		IsEmulator: emulator,
		IsDebug:    debug,
		Network:    network,
		ClientIp:   clientIP,
		DeviceMac:  strings.Join(md["device_mac"], ""),
		Latitude:   strings.Join(md["latitude"], ""),
		Longitude:  strings.Join(md["longitude"], ""),
		Locale:     strings.Join(md["locale"], ""),
		Channel:    strings.Join(md["channel"], ""),
		Product:    strings.Join(md["product"], ""),
		TraceId:    strings.Join(md["trace_id"], ""),
		RiskId:     strings.Join(md["risk_id"], ""),
		Uuid:       strings.Join(md["uuid"], ""),
		Udid:       strings.Join(md["udid"], ""),
		AndroidId:  strings.Join(md["android_id"], ""),
	}
}
