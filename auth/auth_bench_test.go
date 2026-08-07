package auth

import (
	"context"
	"strconv"
	"testing"
	"time"

	md "github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/plugin/credential"
	plugintoken "github.com/appootb/substratum/v2/plugin/token"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/token"
	"github.com/appootb/substratum/v2/util/hash"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setupAuthBench(b *testing.B) (authz *AlgorithmAuth, anonymousCtx, authedCtx context.Context, anonymousMethod, authedMethod string) {
	b.Helper()
	credential.Init()
	plugintoken.Init()

	j := token.Implementor()
	authz = NewAlgorithmAuth(j, j).(*AlgorithmAuth)

	anonymousMethod = "/bench.Service/Anonymous"
	authedMethod = "/bench.Service/Authed"
	authz.RegisterServiceSubjects("bench", map[string][]permission.Subject{
		anonymousMethod: {permission.Subject_NONE},
		authedMethod:    {permission.Subject_SERVER},
	}, nil)

	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "bench",
		Account:   42,
		KeyId:     hash.Sum("auth-bench"),
		Subject:   permission.Subject_SERVER,
		Roles:     []string{"member"},
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		b.Fatal(err)
	}

	basePairs := []string{
		md.KeyProduct, "bench",
		md.KeyPlatform, strconv.Itoa(int(common.Platform_PLATFORM_SERVER)),
		md.KeyLocale, "en",
	}
	anonymousMD := metadata.Pairs(basePairs...)
	authedMD := metadata.Pairs(append(basePairs, md.KeyToken, raw)...)

	// Simulate metadata interceptor: parse and stash on context.
	anonymousIncoming := metadata.NewIncomingContext(context.Background(), anonymousMD)
	authedIncoming := metadata.NewIncomingContext(context.Background(), authedMD)
	anonymousCtx = withParsedMetadata(anonymousIncoming)
	authedCtx = withParsedMetadata(authedIncoming)
	return authz, anonymousCtx, authedCtx, anonymousMethod, authedMethod
}

func withParsedMetadata(ctx context.Context) context.Context {
	// metadata.UnaryServerInterceptor stores under unexported metadataKey;
	// call ParseIncomingMetadata path by invoking the interceptor once.
	out, err := md.UnaryServerInterceptor()(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		return ctx, nil
	})
	if err != nil {
		panic(err)
	}
	return out.(context.Context)
}

func BenchmarkAuthenticate_Anonymous(b *testing.B) {
	authz, ctx, _, method, _ := setupAuthBench(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := authz.Authenticate(ctx, method); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAuthenticate_JWT(b *testing.B) {
	authz, _, ctx, _, method := setupAuthBench(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := authz.Authenticate(ctx, method); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTokenParse(b *testing.B) {
	credential.Init()
	plugintoken.Init()
	j := token.Implementor()
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "bench",
		Account:   42,
		KeyId:     hash.Sum("auth-bench-parse"),
		Subject:   permission.Subject_SERVER,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		b.Fatal(err)
	}
	meta := &common.Metadata{
		Token:    raw,
		Platform: common.Platform_PLATFORM_SERVER,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = j.Parse(meta); err != nil {
			b.Fatal(err)
		}
	}
}
