package storage

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// IsEmpty reports whether err represents a missing cache or DB record
// (redis.Nil or gorm.ErrRecordNotFound), including wrapped errors.
func IsEmpty(err error) bool {
	if err == redis.Nil || err == gorm.ErrRecordNotFound {
		return true
	}
	// Redis
	if errors.Is(err, redis.Nil) {
		return true
	}
	// SQL
	return errors.Is(err, gorm.ErrRecordNotFound)
}
