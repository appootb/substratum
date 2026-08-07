package metadata

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/grpc/metadata"
)

func fullIncomingContext() context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		KeyProduct, "app",
		KeyPackage, "com.example",
		KeyVersion, "1.2.3",
		KeyOSVersion, "17.0",
		KeyBrand, "Apple",
		KeyModel, "iPhone",
		KeyDeviceID, "udid-1",
		KeyFingerprint, "fp-1",
		KeyLocale, "zh-CN",
		KeyLatitude, "31.2",
		KeyLongitude, "121.5",
		KeyPlatform, fmt.Sprintf("%d", int(common.Platform_PLATFORM_MOBILE)),
		KeyNetwork, "2",
		KeyTimestamp, fmt.Sprintf("%d", time.Now().Unix()),
		KeyTraceID, "trace-1",
		KeyIsEmulator, "false",
		KeyChannel, "appstore",
		KeyUUID, "uuid-1",
		KeyIMEI, "imei-1",
		KeyDeviceMac, "aa:bb:cc",
		KeyUserAgent, "ua-custom",
		KeyOriginalIP, "1.2.3.4",
		KeyToken, "jwt-token-value",
	))
}

func BenchmarkParseIncomingMetadata_Full(b *testing.B) {
	ctx := fullIncomingContext()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ParseIncomingMetadata(ctx) == nil {
			b.Fatal("nil metadata")
		}
	}
}

func BenchmarkParseIncomingMetadata_Minimal(b *testing.B) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		KeyPlatform, "1",
		KeyToken, "t",
	))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ParseIncomingMetadata(ctx) == nil {
			b.Fatal("nil metadata")
		}
	}
}

func BenchmarkParseIncomingMetadata_BearerFallback(b *testing.B) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		KeyAuthorization, "Bearer jwt-token-value",
		KeyPlatform, "1",
	))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := ParseIncomingMetadata(ctx)
		if got == nil || got.GetToken() != "jwt-token-value" {
			b.Fatalf("token=%v", got)
		}
	}
}

func BenchmarkParseIncomingMetadata_Missing(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ParseIncomingMetadata(ctx) != nil {
			b.Fatal("expected nil")
		}
	}
}

func BenchmarkLazyAttachOnly(b *testing.B) {
	base := fullIncomingContext()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := ContextWithIncomingMetadata(base)
		if ctx.Value(metadataKey{}) == nil {
			b.Fatal("missing holder")
		}
	}
}

func BenchmarkLazyIncomingMetadata_FirstRead(b *testing.B) {
	base := fullIncomingContext()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := ContextWithIncomingMetadata(base)
		if IncomingMetadata(ctx) == nil {
			b.Fatal("nil")
		}
	}
}

func BenchmarkLazyIncomingMetadata_CachedRead(b *testing.B) {
	ctx := ContextWithIncomingMetadata(fullIncomingContext())
	if IncomingMetadata(ctx) == nil {
		b.Fatal("nil")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if IncomingMetadata(ctx) == nil {
			b.Fatal("nil")
		}
	}
}
