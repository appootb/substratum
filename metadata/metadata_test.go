package metadata

import (
	"context"
	"testing"

	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc/metadata"
)

func TestParseIncomingMetadataBasic(t *testing.T) {
	md := metadata.Pairs(
		KeyProduct, "app",
		KeyPackage, "com.example",
		KeyVersion, "1.2.3",
		KeyPlatform, "1",
		KeyNetwork, "2",
		KeyTimestamp, "1700000000",
		KeyIsEmulator, "true",
		KeyUserAgent, "ua-custom",
		KeyOriginalIP, "1.2.3.4",
		KeyLocale, "zh-CN",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got := ParseIncomingMetadata(ctx)
	if got == nil {
		t.Fatal("nil metadata")
	}
	if got.Product != "app" || got.Package != "com.example" || got.Version != "1.2.3" {
		t.Fatalf("basic fields: %+v", got)
	}
	if got.Platform != common.Platform(1) {
		t.Fatalf("platform=%v", got.Platform)
	}
	if got.Network != common.Network(2) {
		t.Fatalf("network=%v", got.Network)
	}
	if got.Timestamp != 1700000000 || !got.IsEmulator {
		t.Fatalf("ts/emulator: %+v", got)
	}
	if got.UserAgent != "ua-custom" || got.ClientIp != "1.2.3.4" || got.Locale != "zh-CN" {
		t.Fatalf("ua/ip/locale: %+v", got)
	}
}

func TestParseIncomingMetadataMissing(t *testing.T) {
	if ParseIncomingMetadata(context.Background()) != nil {
		t.Fatal("expected nil without incoming metadata")
	}
}

func TestParseIncomingMetadataPlatformEnumName(t *testing.T) {
	// Use a known enum name if present; fall back to numeric already covered.
	md := metadata.Pairs(KeyPlatform, "PLATFORM_UNSPECIFIED")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got := ParseIncomingMetadata(ctx)
	if got == nil {
		t.Fatal("nil")
	}
	if got.Platform != common.Platform_PLATFORM_UNSPECIFIED {
		t.Fatalf("platform=%v", got.Platform)
	}
}

func TestParseIncomingMetadataUserAgentFallback(t *testing.T) {
	md := metadata.Pairs(KeyIANAUserAgent, "mozilla")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got := ParseIncomingMetadata(ctx)
	if got.UserAgent != "mozilla" {
		t.Fatalf("ua=%q", got.UserAgent)
	}
}

func TestParseIncomingMetadataAuthorizationBearer(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"Authorization", "Bearer jwt-token-value",
	))
	parsed := ParseIncomingMetadata(ctx)
	if parsed.GetToken() != "jwt-token-value" {
		t.Fatalf("expected bearer token, got %q", parsed.GetToken())
	}
}
