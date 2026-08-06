package storage

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisOption customizes go-redis client Options during InitRedis.
type RedisOption func(*redis.Options)

// WithRedisMaxRetry sets Options.MaxRetries.
func WithRedisMaxRetry(retries int) RedisOption {
	return func(opt *redis.Options) {
		opt.MaxRetries = retries
	}
}

// WithRedisMaxConnAge sets Options.MaxConnAge (connection max lifetime).
func WithRedisMaxConnAge(dur time.Duration) RedisOption {
	return func(opt *redis.Options) {
		opt.MaxConnAge = dur
	}
}

// WithRedisIdleCheckFreq sets Options.IdleCheckFrequency for the connection reaper.
func WithRedisIdleCheckFreq(dur time.Duration) RedisOption {
	return func(opt *redis.Options) {
		opt.IdleCheckFrequency = dur
	}
}
