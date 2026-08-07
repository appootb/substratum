package gateway

import (
	"context"
	"net/http/httptest"
	"testing"

	md "github.com/appootb/substratum/v2/metadata"
)

func TestURLQueryMetadataIgnoresToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/example?token=secret&plf=7", nil)
	got := URLQueryMetadata(context.Background(), req)
	if _, ok := got[md.KeyToken]; ok {
		t.Fatal("REST query metadata must not include token")
	}
	if v := got[md.KeyPlatform]; len(v) != 1 || v[0] != "7" {
		t.Fatalf("expected platform query param, got %v", v)
	}
}

func TestWebsocketQueryMetadataIncludesToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/stream?token=secret&plf=7", nil)
	got := WebsocketQueryMetadata(context.Background(), req)
	if v := got[md.KeyToken]; len(v) != 1 || v[0] != "secret" {
		t.Fatalf("expected websocket token, got %v", v)
	}
}
