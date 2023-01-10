package storage

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/util/hash"
	"github.com/appootb/substratum/util/ssh"
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

type Storage struct {
	mu sync.RWMutex

	// ElasticSearch
	elastic6 *es6.Client
	elastic7 *es7.Client
	// SQL DBs
	masterDB *gorm.DB
	slaveIdx uint64
	slaveDBs []*gorm.DB
	// Redis caches
	caches []redis.Cmdable
}

func (s *Storage) InitDB(configs []storage.Config, opts ...storage.SQLOption) error {
	for i, cfg := range configs {
		dialect, err := cfg.Dialect()
		if err != nil {
			return err
		}
		//
		db, err := gorm.Open(string(dialect.Type()), dialect.URL())
		if err != nil {
			return err
		}
		for _, o := range opts {
			o(db)
		}
		s.mu.Lock()
		if i > 0 {
			s.slaveDBs = append(s.slaveDBs, db)
		} else {
			s.masterDB = db
		}
		s.mu.Unlock()
	}
	return nil
}

func (s *Storage) InitElasticSearch(config storage.Config, opts ...storage.ElasticOption) error {
	dialect, err := config.Dialect()
	if err != nil {
		return err
	}
	//
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

func (s *Storage) InitRedis(configs []storage.Config, opts ...storage.RedisOption) error {
	if metadata.SSHTunnel != "" {
		opts = append(opts, func(opt *redis.Options) {
			opt.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := ssh.NewTunnel(metadata.SSHTunnel)
				return dialer.Dial(network, addr)
			}
		})
	}
	//
	for _, cfg := range configs {
		dialect, err := cfg.Dialect()
		if err != nil {
			return err
		}
		//
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

func (s *Storage) GetDB(readOnly ...bool) *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(readOnly) == 0 || !readOnly[0] || s.slaveDBs == nil {
		return s.masterDB
	}
	slaves := uint64(len(s.slaveDBs))
	if slaves == 0 {
		return s.masterDB
	}
	idx := atomic.AddUint64(&s.slaveIdx, 1)
	return s.slaveDBs[idx%slaves]
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
