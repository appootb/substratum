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
	InitDB(dialects []Dialect, opts ...SQLOption) error
	InitRedis(dialects []Dialect, opts ...RedisOption) error
	AddDB(db *gorm.DB)
	AddRedis(cache redis.Cmdable)
	GetDBs() []*gorm.DB
	GetDB(key interface{}) *gorm.DB
	GetRedisz() []redis.Cmdable
	GetRedis(key interface{}) redis.Cmdable
}

func newStorage() Storage {
	return &storage{}
}

type storage struct {
	mu sync.RWMutex

	dbs    []*gorm.DB
	caches []redis.Cmdable
}

func (s *storage) InitDB(dialects []Dialect, opts ...SQLOption) error {
	for _, dialect := range dialects {
		db, err := gorm.Open(string(dialect.Type()), dialect.URL())
		if err != nil {
			return err
		}
		for _, opt := range opts {
			opt(db)
		}
		s.AddDB(db)
	}
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
		s.AddRedis(redis.NewClient(options))
	}
	return nil
}

func (s *storage) AddDB(db *gorm.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dbs = append(s.dbs, db)
}

func (s *storage) AddRedis(cache redis.Cmdable) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.caches = append(s.caches, cache)
}
func (s *storage) GetDBs() []*gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dbs
}

func (s *storage) GetDB(key interface{}) *gorm.DB {
	dbs := s.GetDBs()
	sum := hash.Sum(key)
	return dbs[sum%int64(len(dbs))]
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
