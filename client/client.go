package client

import (
	"context"
	"strconv"
	"time"

	ictx "github.com/appootb/substratum/internal/context"
	md "github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/proto/go/common"
	"github.com/appootb/substratum/proto/go/permission"
	"github.com/appootb/substratum/proto/go/secret"
	"github.com/appootb/substratum/service"
	"github.com/appootb/substratum/token"
	"github.com/appootb/substratum/util/iphelper"
	"github.com/appootb/substratum/util/random"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultIssuer = "appootb"
)

func WithContext(ctx context.Context, keyID int64, product ...string) context.Context {
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
	val, _ := token.Implementor().Generate(secretInfo)
	if len(product) > 0 && product[0] != "" {
		outgoingMD[md.KeyProduct] = []string{product[0]}
	}
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
	val, _ := token.Implementor().Generate(secretInfo)
	traceID := incomingMD.GetTraceId()
	if traceID == "" {
		traceID = random.String(32)
	}
	outgoingMD := metadata.New(map[string]string{
		md.KeyToken:      val,
		md.KeyPackage:    incomingMD.GetPackage(),
		md.KeyVersion:    incomingMD.GetVersion(),
		md.KeyOSVersion:  incomingMD.GetOsVersion(),
		md.KeyBrand:      incomingMD.GetBrand(),
		md.KeyModel:      incomingMD.GetModel(),
		md.KeyDeviceID:   incomingMD.GetDeviceId(),
		md.KeyPlatform:   strconv.Itoa(int(platform)),
		md.KeyTimestamp:  strconv.FormatInt(now.UnixNano()/1e6, 10),
		md.KeyIsEmulator: strconv.FormatBool(incomingMD.GetIsEmulator()),
		md.KeyNetwork:    incomingMD.GetNetwork().String(),
		md.KeyLatitude:   incomingMD.GetLatitude(),
		md.KeyLongitude:  incomingMD.GetLongitude(),
		md.KeyLocale:     incomingMD.GetLocale(),
		md.KeyChannel:    incomingMD.GetChannel(),
		md.KeyProduct:    incomingMD.GetProduct(),
		md.KeyTraceID:    traceID,
		md.KeyRiskID:     incomingMD.GetRiskId(),
		md.KeyUUID:       incomingMD.GetUuid(),
		md.KeyUDID:       incomingMD.GetUdid(),
		md.KeyIsDebug:    strconv.FormatBool(incomingMD.GetIsDebug()),
	})
	outgoingMD.Set(md.KeyOriginalIP, incomingMD.GetClientIp(), iphelper.LocalIP())
	return metadata.NewOutgoingContext(ictx.Context, outgoingMD)
}
