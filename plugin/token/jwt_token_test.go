package token

import (
	"testing"
	"time"

	"github.com/appootb/substratum/v2/plugin/credential"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/util/hash"
	"github.com/gbrlsnchs/jwt/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	credential.Init()
}

func TestJwtToken_Generate(t *testing.T) {
	j := &JwtToken{}
	now := time.Now()
	token, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestJwtToken_Generate"),
		Subject:   permission.Subject_MOBILE,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := j.ParseRaw(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestJwtToken_Expire(t *testing.T) {
	j := &JwtToken{}
	now := time.Now()
	token, err := j.Generate(&secret.Info{
		Type:      secret.Type_SERVER,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestJwtToken_Expire"),
		Subject:   permission.Subject_SERVER,
		IssuedAt:  timestamppb.New(now.Add(-time.Hour)),
		ExpiredAt: timestamppb.New(now.Add(-time.Minute)),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.ParseRaw(token)
	if err != jwt.ErrExpValidation {
		t.Fatal(s, err)
	}
}

func TestJwtToken_Before(t *testing.T) {
	j := &JwtToken{}
	now := time.Now()
	token, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestJwtToken_Before"),
		Subject:   permission.Subject_PC,
		IssuedAt:  timestamppb.New(now.Add(time.Minute * 2)),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.ParseRaw(token)
	if err != jwt.ErrNbfValidation {
		t.Fatal(s, err)
	}
}

func TestJwtToken_TokenNotBefore(t *testing.T) {
	j := &JwtToken{}
	now := time.Now()
	token, err := j.Generate(&secret.Info{
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   123456789,
		KeyId:     hash.Sum("TestJwtToken_TokenNotBefore"),
		Subject:   permission.Subject_GUEST,
		IssuedAt:  timestamppb.New(now.Add(time.Minute * 2)),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = j.ParseRaw(token)
	if err == nil {
		t.Fatal("fail: invalid token")
	}
}
