package auth

import (
	"testing"
	"time"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/substratum/util/datetime"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestJwtToken_Generate(t *testing.T) {
	j := NewJwtToken()
	now := time.Now()
	token, err := j.Generate(&permission.Secret{
		Issuer:    proto.String("appootb"),
		Subject:   proto.String("unit_test"),
		AccountId: proto.Uint64(123456789),
		KeyId:     proto.String("TestJwtToken_Generate"),
		Roles:     []string{permission.TokenLevel_LOW_TOKEN.String()},
		Metadata:  map[string]string{},
		IssuedAt:  datetime.WithTime(now).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)

	s, err := j.Verify(token, permission.TokenLevel_LOW_TOKEN)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestJwtToken_Expire(t *testing.T) {
	j := NewJwtToken()
	now := time.Now()
	token, err := j.Generate(&permission.Secret{
		Issuer:    proto.String("appootb"),
		Subject:   proto.String("unit_test"),
		AccountId: proto.Uint64(123456789),
		KeyId:     proto.String("TestJwtToken_Expire"),
		Roles:     []string{permission.TokenLevel_LOW_TOKEN.String()},
		Metadata:  map[string]string{},
		IssuedAt:  datetime.WithTime(now.Add(-time.Hour)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(-time.Minute)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.Verify(token, permission.TokenLevel_NONE_TOKEN)
	if err != jwt.ErrExpValidation {
		t.Fatal(s, err)
	}
}

func TestJwtToken_Before(t *testing.T) {
	j := NewJwtToken()
	now := time.Now()
	token, err := j.Generate(&permission.Secret{
		Issuer:    proto.String("appootb"),
		Subject:   proto.String("unit_test"),
		AccountId: proto.Uint64(123456789),
		KeyId:     proto.String("TestJwtToken_Before"),
		Roles:     []string{permission.TokenLevel_LOW_TOKEN.String()},
		Metadata:  map[string]string{},
		IssuedAt:  datetime.WithTime(now.Add(time.Minute * 2)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := j.Verify(token, permission.TokenLevel_NONE_TOKEN)
	if err != jwt.ErrNbfValidation {
		t.Fatal(s, err)
	}
}

func TestJwtToken_TokenLevel(t *testing.T) {
	j := NewJwtToken()
	now := time.Now()
	token, err := j.Generate(&permission.Secret{
		Issuer:    proto.String("appootb"),
		Subject:   proto.String("unit_test"),
		AccountId: proto.Uint64(123456789),
		KeyId:     proto.String("TestJwtToken_TokenLevel"),
		Roles:     []string{permission.TokenLevel_LOW_TOKEN.String()},
		Metadata:  map[string]string{},
		IssuedAt:  datetime.WithTime(now.Add(time.Minute * 2)).Proto(),
		ExpiredAt: datetime.WithTime(now.Add(time.Hour)).Proto(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = j.Verify(token, permission.TokenLevel_HIGH_TOKEN)
	sErr, ok := status.FromError(err)
	if !ok || sErr.Code() != codes.PermissionDenied {
		t.Fatal(err)
	}
}
