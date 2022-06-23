package client

import (
	"context"
	"strconv"
	"time"

	md "github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/token"
	"github.com/appootb/substratum/v2/util/iphelper"
	"github.com/appootb/substratum/v2/util/random"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultIssuer = "appootb"
)

func WithContext(ctx context.Context, keyID int64) context.Context {
	now := time.Now()
	outgoingMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if reqMD := md.IncomingMetadata(ctx); reqMD != nil {
			return WithMetadata(reqMD, keyID)
		}
		outgoingMD = metadata.MD{}
	}
	account := uint64(0)
	subject := permission.Subject_SERVER
	if accountSecret := service.AccountSecretFromContext(ctx); accountSecret != nil {
		account = accountSecret.GetAccount()
		subject |= accountSecret.GetSubject()
	}
	platform := common.Platform_PLATFORM_SERVER
	if pf := outgoingMD.Get(md.KeyPlatform); len(pf) > 0 {
		if i, err := strconv.Atoi(pf[0]); err != nil {
			platform |= common.Platform(common.Platform_value[pf[0]])
		} else {
			platform |= common.Platform(i)
		}
	}
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    DefaultIssuer,
		Account:   account,
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   subject,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Minute)),
	}
	if pkg := outgoingMD.Get(md.KeyPackage); len(pkg) > 0 {
		secretInfo.Issuer = pkg[0]
	}
	// TODO: cache outgoing tokens
	val, _ := token.Implementor().Generate(secretInfo)
	outgoingMD[md.KeyToken] = []string{val}
	outgoingMD[md.KeyPlatform] = []string{strconv.Itoa(int(platform))}
	outgoingMD[md.KeyTimestamp] = []string{strconv.FormatInt(now.UnixNano()/1e6, 10)}
	outgoingMD[md.KeyOriginalIP] = append(outgoingMD.Get(md.KeyOriginalIP), iphelper.LocalIP())
	return metadata.NewOutgoingContext(ctx, outgoingMD)
}

func WithMetadata(incomingMD *common.Metadata, keyID int64) context.Context {
	now := time.Now()
	platform := common.Platform_PLATFORM_SERVER
	if incomingMD.Platform != common.Platform_PLATFORM_UNSPECIFIED {
		platform |= incomingMD.GetPlatform()
	}
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    DefaultIssuer,
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   permission.Subject_SERVER,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Minute)),
	}
	if incomingMD.Package != "" {
		secretInfo.Issuer = incomingMD.GetPackage()
	}
	// TODO: cache outgoing tokens
	val, _ := token.Implementor().Generate(secretInfo)
	traceID := incomingMD.GetTraceId()
	if traceID == "" {
		traceID = random.String(32)
	}
	outgoingMD := metadata.New(map[string]string{
		md.KeyProduct:     incomingMD.GetProduct(),
		md.KeyPackage:     incomingMD.GetPackage(),
		md.KeyVersion:     incomingMD.GetVersion(),
		md.KeyOSVersion:   incomingMD.GetOsVersion(),
		md.KeyBrand:       incomingMD.GetBrand(),
		md.KeyModel:       incomingMD.GetModel(),
		md.KeyDeviceID:    incomingMD.GetDeviceId(),
		md.KeyFingerprint: incomingMD.GetFingerprint(),
		md.KeyLocale:      incomingMD.GetLocale(),
		md.KeyLatitude:    incomingMD.GetLatitude(),
		md.KeyLongitude:   incomingMD.GetLongitude(),
		md.KeyPlatform:    strconv.Itoa(int(platform)),
		md.KeyNetwork:     incomingMD.GetNetwork().String(),
		md.KeyTimestamp:   strconv.FormatInt(now.UnixNano()/1e6, 10),
		md.KeyTraceID:     traceID,
		md.KeyIsEmulator:  strconv.FormatBool(incomingMD.GetIsEmulator()),
		md.KeyIsDevelop:   strconv.FormatBool(incomingMD.GetIsDevelop()),
		md.KeyIsTesting:   strconv.FormatBool(incomingMD.GetIsTesting()),
		md.KeyChannel:     incomingMD.GetChannel(),
		md.KeyUUID:        incomingMD.GetUuid(),
		md.KeyIMEI:        incomingMD.GetImei(),
		md.KeyDeviceMac:   incomingMD.GetDeviceMac(),
		md.KeyUserAgent:   incomingMD.GetUserAgent(),
		md.KeyToken:       val,
	})
	outgoingMD.Set(md.KeyOriginalIP, incomingMD.GetClientIp(), iphelper.LocalIP())
	return metadata.NewOutgoingContext(context.Background(), outgoingMD)
}
