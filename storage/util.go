package storage

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

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
