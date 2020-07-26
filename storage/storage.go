package storage

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Storage interface.
type Storage interface {
	InitDB(dialect Dialect, opts ...SQLOption) error
	InitRedis(dialects []Dialect, opts ...RedisOption) error
	GetDB() *gorm.DB
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}
