package storage

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisOptions(t *testing.T) {
	opt := &redis.Options{}
	WithRedisMaxRetry(5)(opt)
	WithRedisMaxConnAge(time.Hour)(opt)
	WithRedisIdleCheckFreq(time.Minute)(opt)

	if opt.MaxRetries != 5 {
		t.Fatalf("MaxRetries=%d", opt.MaxRetries)
	}
	if opt.MaxConnAge != time.Hour {
		t.Fatalf("MaxConnAge=%v", opt.MaxConnAge)
	}
	if opt.IdleCheckFrequency != time.Minute {
		t.Fatalf("IdleCheckFrequency=%v", opt.IdleCheckFrequency)
	}
}
