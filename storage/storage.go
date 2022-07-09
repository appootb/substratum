package storage

import (
	"github.com/appootb/substratum/v2/configure"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Storage interface.
type Storage interface {
	InitDB(master configure.Address, slaves []configure.Address, opts ...SQLOption) error
	InitRedis(configs []configure.Address, opts ...RedisOption) error
	GetDB(readOnly ...bool) *gorm.DB
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}
