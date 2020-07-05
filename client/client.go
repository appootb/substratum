package client

import (
	"strconv"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/util/datetime"
	"github.com/appootb/substratum/util/iphelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CallOption(md *common.Metadata, keyID int64) grpc.CallOption {
	now := time.Now()
	secretInfo := &secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   md.GetAccount(),
		KeyId:     keyID,
		Roles:     []string{},
		Subject:   permission.Subject_SERVER,
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Minute)).Proto(),
	}
	token, _ := auth.DefaultToken.Generate(secretInfo)
	m := map[string]string{
		"account":         strconv.FormatUint(md.GetAccount(), 10),
		"token":           token,
		"package":         md.GetPackage(),
		"version":         md.GetVersion(),
		"device_id":       md.GetDeviceId(),
		"x-forwarded-for": iphelper.LocalIP(),
		"platform":        common.Platform_PLATFORM_SERVER.String(),
		"timestamp":       strconv.FormatInt(now.UnixNano()/1e6, 10),
		"locale":          md.GetLocale(),
		"channel":         md.GetChannel(),
		"product":         md.GetProduct(),
		"trace_id":        md.GetTraceId(),
	}
	header := metadata.New(m)
	return grpc.Header(&header)
}
