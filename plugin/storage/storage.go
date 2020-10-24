package storage

import (
	"sync"

	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/util/hash"
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

type Storage struct {
	mu sync.RWMutex

	elastic6 *es6.Client
	elastic7 *es7.Client
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

func (s *Storage) InitElasticSearch(dialect storage.Dialect, opts ...storage.ElasticOption) error {
	switch dialect.Type() {
	case storage.DialectElasticSearch6:
		cfg6 := es6.Config{
			Addresses: []string{dialect.URL()},
			Username:  dialect.Meta().Username,
			Password:  dialect.Meta().Password,
		}
		for _, o := range opts {
			o(&cfg6, nil)
		}
		cli6, err := es6.NewClient(cfg6)
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.elastic6 = cli6
		s.mu.Unlock()
	case storage.DialectElasticSearch7:
		cfg7 := es7.Config{
			Addresses: []string{dialect.URL()},
			Username:  dialect.Meta().Username,
			Password:  dialect.Meta().Password,
		}
		for _, o := range opts {
			o(nil, &cfg7)
		}
		cli7, err := es7.NewClient(cfg7)
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.elastic7 = cli7
		s.mu.Unlock()
	}
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

func (s *Storage) GetESv6() *es6.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.elastic6
}

func (s *Storage) GetESv7() *es7.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.elastic7
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
