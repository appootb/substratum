package storage

import (
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Storage interface.
type Storage interface {
	InitDB(master Config, slaves []Config, opts ...SQLOption) error
	InitElasticSearch(config Config, opts ...ElasticOption) error
	InitRedis(configs []Config, opts ...RedisOption) error
	GetDB(readOnly ...bool) *gorm.DB
	GetESv6() *es6.Client
	GetESv7() *es7.Client
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}
