package logger

import "sync/atomic"

// accessLogEnabled is 1 (default) when unary/stream access logging interceptors
// perform work; 0 skips logging entirely (handler is invoked directly).
var accessLogEnabled uint32 = 1

// SetAccessLogEnabled enables or disables gRPC access/stream logging interceptors.
// Default is true. Disabling is useful for load tests or high-QPS services that
// collect logs elsewhere. Application Logger.Info/Error calls are unaffected.
func SetAccessLogEnabled(enabled bool) {
	var v uint32
	if enabled {
		v = 1
	}
	atomic.StoreUint32(&accessLogEnabled, v)
}

// AccessLogEnabled reports whether access-log interceptors are active.
func AccessLogEnabled() bool {
	return atomic.LoadUint32(&accessLogEnabled) != 0
}
