package storage

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

func IsEmpty(err error) bool {
	// Redis
	if errors.Is(err, redis.Nil) {
		return true
	}
	// SQL
	if err == gorm.ErrRecordNotFound {
		return true
	}
	var gormErr *gorm.Errors
	if errors.As(err, &gormErr) {
		return gorm.IsRecordNotFoundError(gormErr)
	}
	return false
}
