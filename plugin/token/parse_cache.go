package token

import (
	"sync/atomic"
	"time"

	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/util/cache"
)

const (
	// DefaultParseCacheSize is the max number of verified tokens retained.
	DefaultParseCacheSize = 4096
	// DefaultParseCacheMaxTTL caps how long a verified token may be reused
	// without re-verifying (revocation delay window).
	DefaultParseCacheMaxTTL = 10 * time.Second
	// DefaultParseCacheSkew is subtracted from JWT exp when computing entry TTL.
	DefaultParseCacheSkew = 2 * time.Second
)

var (
	parseCache = cache.New(cache.LRU, cache.WithSize(DefaultParseCacheSize))

	// 1 = enabled (default).
	parseCacheEnabled = uint32(1)
	// nanoseconds; defaults match DefaultParseCacheMaxTTL / Skew.
	parseCacheMaxTTLNs = int64(DefaultParseCacheMaxTTL)
	parseCacheSkewNs   = int64(DefaultParseCacheSkew)
)

// SetParseCacheEnabled turns process-local JWT verification caching on or off.
// Default is true. Cache hits skip signature verification for the same token
// string until the entry TTL expires.
func SetParseCacheEnabled(enabled bool) {
	var v uint32
	if enabled {
		v = 1
	}
	atomic.StoreUint32(&parseCacheEnabled, v)
	if !enabled {
		parseCache.Purge()
	}
}

// ParseCacheEnabled reports whether ParseRaw verification caching is active.
func ParseCacheEnabled() bool {
	return atomic.LoadUint32(&parseCacheEnabled) != 0
}

// SetParseCacheMaxTTL sets the maximum TTL for a cached verification result.
// Non-positive values reset to DefaultParseCacheMaxTTL.
// Effective TTL is min(maxTTL, time.Until(exp)-skew).
func SetParseCacheMaxTTL(d time.Duration) {
	if d <= 0 {
		d = DefaultParseCacheMaxTTL
	}
	atomic.StoreInt64(&parseCacheMaxTTLNs, int64(d))
}

// ParseCacheMaxTTL returns the configured maximum cache TTL.
func ParseCacheMaxTTL() time.Duration {
	return time.Duration(atomic.LoadInt64(&parseCacheMaxTTLNs))
}

// SetParseCacheSkew sets the safety margin subtracted from JWT expiry when
// computing cache TTL. Non-positive values reset to DefaultParseCacheSkew.
func SetParseCacheSkew(d time.Duration) {
	if d <= 0 {
		d = DefaultParseCacheSkew
	}
	atomic.StoreInt64(&parseCacheSkewNs, int64(d))
}

// ParseCacheSkew returns the configured expiry skew.
func ParseCacheSkew() time.Duration {
	return time.Duration(atomic.LoadInt64(&parseCacheSkewNs))
}

// ClearParseCache drops all cached verification results.
func ClearParseCache() {
	parseCache.Purge()
}

func parseCacheTTL(info *secret.Info) time.Duration {
	maxTTL := ParseCacheMaxTTL()
	skew := ParseCacheSkew()
	if info == nil || info.GetExpiredAt() == nil {
		return maxTTL
	}
	ttl := time.Until(info.GetExpiredAt().AsTime()) - skew
	if ttl <= 0 {
		return 0
	}
	if ttl > maxTTL {
		return maxTTL
	}
	return ttl
}

func cachedParseRaw(raw string, load func() (*secret.Info, error)) (*secret.Info, error) {
	if !ParseCacheEnabled() || raw == "" {
		return load()
	}
	// Fast path: hot token without single-flight overhead when already present.
	if v, ok := parseCache.Get(raw); ok {
		info := v.(*secret.Info)
		if info.GetExpiredAt() == nil || time.Now().Before(info.GetExpiredAt().AsTime()) {
			return info, nil
		}
		parseCache.Del(raw)
	}
	v, err := parseCache.GetOrLoad(raw, func(_ interface{}) (interface{}, time.Duration, error) {
		info, err := load()
		if err != nil {
			// Do not cache failures (avoids sticky auth errors / poisoning).
			return nil, 0, err
		}
		return info, parseCacheTTL(info), nil
	})
	if err != nil {
		return nil, err
	}
	info := v.(*secret.Info)
	if info.GetExpiredAt() != nil && !time.Now().Before(info.GetExpiredAt().AsTime()) {
		parseCache.Del(raw)
		return load()
	}
	return info, nil
}
