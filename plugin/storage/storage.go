package storage

import (
	"sync"
	"sync/atomic"

	"github.com/appootb/substratum/v2/configure"
	"github.com/appootb/substratum/v2/errors"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/util/hash"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type Storage struct {
	mu sync.RWMutex

	// SQL DBs
	masterDB *gorm.DB
	slaveIdx uint64
	slaveDBs []*gorm.DB
	// Redis caches
	caches []redis.Cmdable

	// Common instances
	common map[configure.Schema]interface{}
}

func (s *Storage) InitDB(master configure.Address, slaves []configure.Address, opts ...storage.SQLOption) error {
	cfg := &gorm.Config{}
	for _, o := range opts {
		o(cfg, nil)
	}
	//
	var (
		err error
		db  *gorm.DB
	)
	// Master
	db, err = gorm.Open(storage.SQLDialectImplementor().Open(master), cfg)
	if err != nil {
		return err
	}
	for _, o := range opts {
		o(nil, db)
	}
	s.mu.Lock()
	s.masterDB = db
	s.mu.Unlock()
	// Slaves
	if slaves == nil || len(slaves) == 0 {
		return nil
	}
	s.slaveDBs = make([]*gorm.DB, 0, len(slaves))
	for _, slave := range slaves {
		db, err = gorm.Open(storage.SQLDialectImplementor().Open(slave), cfg)
		if err != nil {
			return err
		}
		for _, o := range opts {
			o(nil, db)
		}
		s.mu.Lock()
		s.slaveDBs = append(s.slaveDBs, db)
		s.mu.Unlock()
	}
	return nil
}

func (s *Storage) InitRedis(configs []configure.Address, opts ...storage.RedisOption) error {
	for _, cfg := range configs {
		dialect := redisCache{cfg}
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

func (s *Storage) InitCommon(config configure.Address) error {
	impl := storage.CommonDialectImplementor(config.Schema)
	if impl == nil {
		return errors.Newf(codes.FailedPrecondition, "unknown storage schema: %s", config.Schema)
	}
	dialect, err := impl.Open(config)
	if err != nil {
		return err
	}
	//
	s.mu.Lock()
	s.common[config.Schema] = dialect
	s.mu.Unlock()
	return nil
}

func (s *Storage) GetCommon(schema configure.Schema) interface{} {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.common[schema]
}
