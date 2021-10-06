package storage

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func IsEmpty(err error) bool {
	// Redis
	if errors.Is(err, redis.Nil) {
		return true
	}
	// SQL
	return errors.Is(err, gorm.ErrRecordNotFound)
}
