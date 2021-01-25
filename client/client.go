package client

import (
	"context"
	"strconv"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/protobuf/go/service"
	appootb "github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/token"
	"github.com/appootb/substratum/util/datetime"
	"github.com/appootb/substratum/util/iphelper"
	"github.com/appootb/substratum/util/random"
	"google.golang.org/grpc/metadata"
)

func WithContext(ctx context.Context, keyID int64) context.Context {
	now := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if reqMD := appootb.RequestMetadata(ctx); reqMD != nil {
			return WithMetadata(reqMD, keyID)
		}
		md = metadata.MD{}
	}
	account := uint64(0)
	subject := permission.Subject_SERVER
	if accountSecret := service.AccountSecretFromContext(ctx); accountSecret != nil {
		account = accountSecret.GetAccount()
		subject |= accountSecret.GetSubject()
	}
	platform := common.Platform_PLATFORM_SERVER
	if pf := md.Get("platform"); len(pf) > 0 {
		if i, err := strconv.Atoi(pf[0]); err != nil {
			platform |= common.Platform(common.Platform_value[pf[0]])
		} else {
			platform |= common.Platform(i)
		}
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
		Subject:   subject,
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Minute)).Proto(),
	}
	val, _ := token.Implementor().Generate(secretInfo)
	md["token"] = []string{val}
	md["platform"] = []string{strconv.Itoa(int(platform))}
	md["timestamp"] = []string{strconv.FormatInt(now.UnixNano()/1e6, 10)}
	md["x-forwarded-for"] = append(md.Get("x-forwarded-for"), iphelper.LocalIP())
	return metadata.NewOutgoingContext(ctx, md)
}

func WithMetadata(md *common.Metadata, keyID int64) context.Context {
	now := time.Now()
	issuer := "appootb"
	if md.Package != nil {
		issuer = md.GetPackage()
	}
	platform := common.Platform_PLATFORM_SERVER
	if md.Platform != nil {
		platform |= md.GetPlatform()
	}
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    issuer,
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   permission.Subject_SERVER,
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Minute)).Proto(),
	}
	val, _ := token.Implementor().Generate(secretInfo)
	traceID := md.GetTraceId()
	if traceID == "" {
		traceID = random.String(32)
	}
	outgoingMD := metadata.New(map[string]string{
		"token":       val,
		"package":     md.GetPackage(),
		"version":     md.GetVersion(),
		"os_version":  md.GetOsVersion(),
		"brand":       md.GetBrand(),
		"model":       md.GetModel(),
		"device_id":   md.GetDeviceId(),
		"platform":    strconv.Itoa(int(platform)),
		"timestamp":   strconv.FormatInt(now.UnixNano()/1e6, 10),
		"is_emulator": strconv.FormatBool(md.GetIsEmulator()),
		"network":     md.GetNetwork().String(),
		"latitude":    md.GetLatitude(),
		"longitude":   md.GetLongitude(),
		"locale":      md.GetLocale(),
		"channel":     md.GetChannel(),
		"product":     md.GetProduct(),
		"trace_id":    traceID,
		"is_debug":    strconv.FormatBool(md.GetIsDebug()),
	})
	outgoingMD.Set("x-forwarded-for", md.GetClientIp(), iphelper.LocalIP())
	return metadata.NewOutgoingContext(context.Background(), outgoingMD)
}
