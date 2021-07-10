package storage

import (
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Storage interface.
type Storage interface {
	InitDB(configs []Config, opts ...SQLOption) error
	InitElasticSearch(config Config, opts ...ElasticOption) error
	InitRedis(configs []Config, opts ...RedisOption) error
	GetDB(readOnly ...bool) *gorm.DB
	GetESv6() *es6.Client
	GetESv7() *es7.Client
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}
