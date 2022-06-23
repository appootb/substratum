package storage

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Storage interface.
type Storage interface {
	InitDB(master Config, slaves []Config, opts ...SQLOption) error
	InitRedis(configs []Config, opts ...RedisOption) error
	GetDB(readOnly ...bool) *gorm.DB
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}
