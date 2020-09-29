package storage

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisOption func(*redis.Options)

func WithRedisMaxRetry(retries int) RedisOption {
	return func(opt *redis.Options) {
		opt.MaxRetries = retries
	}
}

func WithRedisMaxConnAge(dur time.Duration) RedisOption {
	return func(opt *redis.Options) {
		opt.MaxConnAge = dur
	}
}

func WithRedisIdleCheckFreq(dur time.Duration) RedisOption {
	return func(opt *redis.Options) {
		opt.IdleCheckFrequency = dur
	}
}
