package cache

import "time"

// entry holds a cached value and optional absolute expiry.
// Shared by LRU/LFU implementations (each may embed additional fields).
type entry struct {
	key    interface{}
	value  interface{}
	expire time.Time
}

// ExpiredAt reports whether the entry is expired at ts.
// Zero expire means never expires. Zero ts uses time.Now().
func (e *entry) ExpiredAt(ts time.Time) bool {
	if e.expire.IsZero() {
		return false
	}
	if ts.IsZero() {
		ts = time.Now()
	}
	return e.expire.Before(ts)
}

// expireAt converts a relative expire duration to an absolute time.
// expire > 0 is a TTL from now; expire < 0 (e.g. DurationPersistence) means never expires.
// expire == 0 is not valid here (callers treat 0 as delete).
func expireAt(expire time.Duration) time.Time {
	if expire > 0 {
		return time.Now().Add(expire)
	}
	return time.Time{}
}
