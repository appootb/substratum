package client

import (
	"context"
	"strconv"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/substratum/token"
	"github.com/appootb/substratum/util/datetime"
	"github.com/appootb/substratum/util/iphelper"
	"google.golang.org/grpc/metadata"
)

func WithContext(ctx context.Context, keyID int64) context.Context {
	now := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	account := uint64(0)
	if id := md.Get("account"); len(id) > 0 {
		account, _ = strconv.ParseUint(id[0], 10, 64)
	}
	issuer := "appootb"
	if pkg := md.Get("package"); len(pkg) > 0 {
		issuer = pkg[0]
	}
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    issuer,
		Account:   account,
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   permission.Subject_SERVER,
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Minute)).Proto(),
	}
	val, _ := token.Implementor().Generate(secretInfo)
	md["token"] = []string{val}
	md["platform"] = []string{common.Platform_PLATFORM_SERVER.String()}
	md["timestamp"] = []string{strconv.FormatInt(now.UnixNano()/1e6, 10)}
	md["x-forwarded-for"] = []string{iphelper.LocalIP()}
	return metadata.NewOutgoingContext(ctx, md)
}

func WithMetadata(md *common.Metadata, keyID int64) metadata.MD {
	now := time.Now()
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    md.GetPackage(),
		Account:   md.GetAccount(),
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   permission.Subject_SERVER,
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Minute)).Proto(),
	}
	val, _ := token.Implementor().Generate(secretInfo)
	m := map[string]string{
		"account":         strconv.FormatUint(md.GetAccount(), 10),
		"token":           val,
		"package":         md.GetPackage(),
		"version":         md.GetVersion(),
		"os_version":      md.GetOsVersion(),
		"brand":           md.GetBrand(),
		"model":           md.GetModel(),
		"device_id":       md.GetDeviceId(),
		"x-forwarded-for": iphelper.LocalIP(),
		"platform":        common.Platform_PLATFORM_SERVER.String(),
		"timestamp":       strconv.FormatInt(now.UnixNano()/1e6, 10),
		"is_emulator":     strconv.FormatBool(md.GetIsEmulator()),
		"network":         md.GetNetwork().String(),
		"client_ip":       md.GetClientIp(),
		"latitude":        md.GetLatitude(),
		"longitude":       md.GetLongitude(),
		"locale":          md.GetLocale(),
		"channel":         md.GetChannel(),
		"product":         md.GetProduct(),
		"trace_id":        md.GetTraceId(),
		"is_debug":        strconv.FormatBool(md.GetIsDebug()),
	}
	return metadata.New(m)
}
