package monitor

import "sync/atomic"

// metricsEnabled is 1 (default) when Prometheus gRPC interceptors record metrics;
// 0 skips metrics and invokes the handler directly.
var metricsEnabled uint32 = 1

// SetMetricsEnabled enables or disables gRPC server metrics interceptors.
// Default is true. Disabling is useful for load tests or services that do not
// scrape Prometheus. InitializeMetrics is a no-op while metrics are disabled.
func SetMetricsEnabled(enabled bool) {
	var v uint32
	if enabled {
		v = 1
	}
	atomic.StoreUint32(&metricsEnabled, v)
}

// MetricsEnabled reports whether metrics interceptors are active.
func MetricsEnabled() bool {
	return atomic.LoadUint32(&metricsEnabled) != 0
}
