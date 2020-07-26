package storage

import (
	"sync"

	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/util/hash"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type Storage struct {
	mu sync.RWMutex

	database *gorm.DB
	caches   []redis.Cmdable
}

func (s *Storage) InitDB(dialect storage.Dialect, opts ...storage.SQLOption) error {
	db, err := gorm.Open(string(dialect.Type()), dialect.URL())
	if err != nil {
		return err
	}
	for _, o := range opts {
		o(db)
	}
	s.mu.Lock()
	s.database = db
	s.mu.Unlock()
	return nil
}

func (s *Storage) InitRedis(dialects []storage.Dialect, opts ...storage.RedisOption) error {
	for _, dialect := range dialects {
		options, err := redis.ParseURL(dialect.URL())
		if err != nil {
			return err
		}
		for _, o := range opts {
			o(options)
		}
		s.mu.Lock()
		s.caches = append(s.caches, redis.NewClient(options))
		s.mu.Unlock()
	}
	return nil
}

func (s *Storage) GetDB() *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.database
}

func (s *Storage) GetRedisz() []redis.Cmdable {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.caches
}

func (s *Storage) GetRedis(key interface{}) redis.Cmdable {
	caches := s.GetRedisz()
	sum := hash.Sum(key)
	return caches[sum%int64(len(caches))]
}
