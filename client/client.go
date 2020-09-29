package client

import (
	"context"
	"strconv"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/protobuf/go/service"
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
	if accountSecret := service.AccountSecretFromContext(ctx); accountSecret != nil {
		account = accountSecret.GetAccount()
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
	md["x-forwarded-for"] = append(md.Get("x-forwarded-for"), iphelper.LocalIP())
	return metadata.NewOutgoingContext(ctx, md)
}
