package storage

import (
	"context"
	"net"
	"time"

	"github.com/appootb/substratum/util/ssh"
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

func WithSSHTunnel(v string) RedisOption {
	return func(opt *redis.Options) {
		opt.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := ssh.NewTunnel(v)
			return dialer.Dial(network, addr)
		}
	}
}
