package logger

import (
	"testing"

	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
)

func TestRedactSecret(t *testing.T) {
	if RedactSecret(nil) != nil {
		t.Fatal("expected nil for missing secret")
	}
	if RedactSecret(&secret.Info{}) != nil {
		t.Fatal("expected nil for anonymous secret")
	}

	redacted := RedactSecret(&secret.Info{
		Type:    secret.Type_CLIENT,
		Account: 42,
		KeyId:   99,
		Issuer:  "appootb",
		Subject: permission.Subject_WEB,
		Roles:   []string{"member"},
	})
	if redacted["account"] != uint64(42) {
		t.Fatalf("unexpected account: %v", redacted["account"])
	}
	if _, ok := redacted["key_id"]; ok {
		t.Fatal("key_id must not be logged")
	}
	if _, ok := redacted["issuer"]; ok {
		t.Fatal("issuer must not be logged")
	}
}
