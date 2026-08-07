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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthenticateServerTokenWithPlatform(t *testing.T) {
	credential.Init()
	plugintoken.Init()

	j := token.Implementor()
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestAuthenticateServerTokenWithPlatform"),
		Subject:   permission.Subject_SERVER,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}

	authz := NewAlgorithmAuth(j, j)
	method := "/example.Service/Internal"
	authz.RegisterServiceSubjects("example", map[string][]permission.Subject{
		method: {permission.Subject_SERVER},
	}, nil)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		md.KeyToken, raw,
		md.KeyPlatform, strconv.Itoa(int(common.Platform_PLATFORM_SERVER)),
	))
	interceptor := md.UnaryServerInterceptor()
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, req interface{}) (interface{}, error) {
		info, err := authz.Authenticate(ctx, method)
		if err != nil {
			return nil, err
		}
		if info.GetType() != secret.Type_SERVER {
			t.Fatalf("expected server token, got %v", info.GetType())
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthenticateRejectsTokenTypeMismatch(t *testing.T) {
	credential.Init()
	plugintoken.Init()

	j := token.Implementor()
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestAuthenticateRejectsTokenTypeMismatch"),
		Subject:   permission.Subject_SERVER,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}

	authz := NewAlgorithmAuth(j, j)
	method := "/example.Service/Ping"
	authz.RegisterServiceSubjects("example", map[string][]permission.Subject{
		method: {permission.Subject_MOBILE},
	}, nil)

	// Client platform with server token must fail at parse-time type check.
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		md.KeyToken, raw,
		md.KeyPlatform, strconv.Itoa(int(common.Platform_PLATFORM_MOBILE)),
	))
	interceptor := md.UnaryServerInterceptor()
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, req interface{}) (interface{}, error) {
		_, err := authz.Authenticate(ctx, method)
		if err == nil {
			t.Fatal("expected token type mismatch")
		}
		return nil, err
	})
	if err == nil {
		t.Fatal("expected authenticate failure")
	}
}
