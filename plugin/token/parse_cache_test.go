package token

import (
	"testing"
	"time"

	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/util/hash"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseRawCacheHit(t *testing.T) {
	ClearParseCache()
	SetParseCacheEnabled(true)
	t.Cleanup(func() {
		SetParseCacheEnabled(true)
		ClearParseCache()
	})

	j := &JwtToken{}
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   42,
		KeyId:     hash.Sum("TestParseRawCacheHit"),
		Subject:   permission.Subject_WEB,
		Roles:     []string{"member"},
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}

	a, err := j.ParseRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	b, err := j.ParseRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	if a.GetAccount() != 42 || b.GetAccount() != 42 {
		t.Fatalf("account a=%d b=%d", a.GetAccount(), b.GetAccount())
	}
	// Same cached pointer on hit (immutable usage contract).
	if a != b {
		t.Fatal("expected cache hit to return same *secret.Info")
	}
	if parseCache.Len() < 1 {
		t.Fatal("expected cache entry")
	}
}

func TestParseRawCacheDisabled(t *testing.T) {
	ClearParseCache()
	SetParseCacheEnabled(false)
	t.Cleanup(func() {
		SetParseCacheEnabled(true)
		ClearParseCache()
	})

	j := &JwtToken{}
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "appootb",
		Account:   7,
		KeyId:     hash.Sum("TestParseRawCacheDisabled"),
		Subject:   permission.Subject_WEB,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := j.ParseRaw(raw); err != nil {
		t.Fatal(err)
	}
	if parseCache.Len() != 0 {
		t.Fatalf("cache should stay empty when disabled, len=%d", parseCache.Len())
	}
}

func TestParseCacheTTLRespectsMaxAndExpiry(t *testing.T) {
	SetParseCacheMaxTTL(10 * time.Second)
	SetParseCacheSkew(2 * time.Second)
	t.Cleanup(func() {
		SetParseCacheMaxTTL(DefaultParseCacheMaxTTL)
		SetParseCacheSkew(DefaultParseCacheSkew)
	})

	now := time.Now()
	info := &secret.Info{
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	}
	ttl := parseCacheTTL(info)
	if ttl != 10*time.Second {
		t.Fatalf("ttl=%v want max 10s", ttl)
	}

	info.ExpiredAt = timestamppb.New(now.Add(5 * time.Second))
	ttl = parseCacheTTL(info)
	// ~5s - 2s skew
	if ttl < 2*time.Second || ttl > 4*time.Second {
		t.Fatalf("ttl=%v want ~3s", ttl)
	}

	info.ExpiredAt = timestamppb.New(now.Add(time.Second))
	if parseCacheTTL(info) != 0 {
		t.Fatal("expected zero ttl when within skew of expiry")
	}
}

func BenchmarkParseRaw_CacheMiss(b *testing.B) {
	ClearParseCache()
	SetParseCacheEnabled(false)
	b.Cleanup(func() {
		SetParseCacheEnabled(true)
		ClearParseCache()
	})
	j := &JwtToken{}
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "bench",
		Account:   1,
		KeyId:     hash.Sum("BenchmarkParseRaw_CacheMiss"),
		Subject:   permission.Subject_WEB,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := j.ParseRaw(raw); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseRaw_CacheHit(b *testing.B) {
	ClearParseCache()
	SetParseCacheEnabled(true)
	b.Cleanup(func() {
		SetParseCacheEnabled(true)
		ClearParseCache()
	})
	j := &JwtToken{}
	now := time.Now()
	raw, err := j.Generate(&secret.Info{
		Type:      secret.Type_CLIENT,
		Algorithm: secret.Algorithm_HMAC,
		Issuer:    "bench",
		Account:   1,
		KeyId:     hash.Sum("BenchmarkParseRaw_CacheHit"),
		Subject:   permission.Subject_WEB,
		IssuedAt:  timestamppb.New(now),
		ExpiredAt: timestamppb.New(now.Add(time.Hour)),
	})
	if err != nil {
		b.Fatal(err)
	}
	if _, err := j.ParseRaw(raw); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := j.ParseRaw(raw); err != nil {
			b.Fatal(err)
		}
	}
}
