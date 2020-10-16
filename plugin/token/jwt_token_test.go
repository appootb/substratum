package token

import (
	"testing"
	"time"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/substratum/plugin/credential"
	"github.com/appootb/substratum/util/datetime"
	"github.com/appootb/substratum/util/hash"
	"github.com/gbrlsnchs/jwt/v3"
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
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := j.Parse(token)
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
		IssuedAt:  datetime.WithTime(now.Add(-time.Hour)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(-time.Minute)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.Parse(token)
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
		IssuedAt:  datetime.WithTime(now.Add(time.Minute * 2)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.Parse(token)
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
		IssuedAt:  datetime.WithTime(now.Add(time.Minute * 2)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = j.Parse(token)
	if err == nil {
		t.Fatal("fail: invalid token")
	}
}
