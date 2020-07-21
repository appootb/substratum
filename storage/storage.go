package storage

import (
	"sync"

	"github.com/appootb/substratum/util/hash"
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

func newStorage() Storage {
	return &storage{}
}

type storage struct {
	mu sync.RWMutex

	database *gorm.DB
	caches   []redis.Cmdable
}

func (s *storage) InitDB(dialect Dialect, opts ...SQLOption) error {
	db, err := gorm.Open(string(dialect.Type()), dialect.URL())
	if err != nil {
		return err
	}
	for _, opt := range opts {
		opt(db)
	}
	s.mu.Lock()
	s.database = db
	s.mu.Unlock()
	return nil
}

func (s *storage) InitRedis(dialects []Dialect, opts ...RedisOption) error {
	for _, dialect := range dialects {
		options, err := redis.ParseURL(dialect.URL())
		if err != nil {
			return err
		}
		for _, opt := range opts {
			opt(options)
		}
		s.mu.Lock()
		s.caches = append(s.caches, redis.NewClient(options))
		s.mu.Unlock()
	}
	return nil
}

func (s *storage) GetDB() *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.database
}

func (s *storage) GetRedisz() []redis.Cmdable {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.caches
}

func (s *storage) GetRedis(key interface{}) redis.Cmdable {
	caches := s.GetRedisz()
	sum := hash.Sum(key)
	return caches[sum%int64(len(caches))]
}
