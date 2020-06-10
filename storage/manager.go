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

const (
	DefaultComponent = "default"
)

var (
	Default = New()
)

type Manager struct {
	m sync.RWMutex

	// Database manager for components, e.g. mysql or postgres.
	dbMgr map[string][]*gorm.DB

	// Redis manager for components.
	redisMgr map[string][]redis.Cmdable
}

func New() *Manager {
	return &Manager{
		dbMgr:    make(map[string][]*gorm.DB),
		redisMgr: make(map[string][]redis.Cmdable),
	}
}

func (mgr *Manager) InitDB(component string, dialects []Dialect, opts ...SQLOption) error {
	for _, dialect := range dialects {
		db, err := gorm.Open(string(dialect.Type()), dialect.URL())
		if err != nil {
			return err
		}
		for _, opt := range opts {
			opt(db)
		}
		mgr.AddDB(component, db)
	}

	return nil
}

func (mgr *Manager) AddDB(component string, db *gorm.DB) {
	mgr.m.Lock()
	mgr.dbMgr[component] = append(mgr.dbMgr[component], db)
	mgr.m.Unlock()
}

func (mgr *Manager) InitRedis(component string, dialects []Dialect, opts ...RedisOption) error {
	for _, dialect := range dialects {
		options, err := redis.ParseURL(dialect.URL())
		if err != nil {
			return err
		}
		for _, opt := range opts {
			opt(options)
		}
		mgr.AddRedis(component, redis.NewClient(options))
	}
	return nil
}

func (mgr *Manager) AddRedis(component string, cache redis.Cmdable) {
	mgr.m.Lock()
	mgr.redisMgr[component] = append(mgr.redisMgr[component], cache)
	mgr.m.Unlock()
}

func (mgr *Manager) GetDBs(component string) []*gorm.DB {
	mgr.m.RLock()
	defer mgr.m.RUnlock()
	if dbs, ok := mgr.dbMgr[component]; ok {
		return dbs
	}
	return mgr.dbMgr[DefaultComponent]
}

func (mgr *Manager) GetDB(component string, key interface{}) *gorm.DB {
	dbs := mgr.GetDBs(component)
	if dbs == nil || len(dbs) == 0 {
		panic("no database configuration for component: " + component)
	}
	sum := hash.Sum(key)
	return dbs[sum%int64(len(dbs))]
}

func (mgr *Manager) GetCaches(component string) []redis.Cmdable {
	mgr.m.RLock()
	defer mgr.m.RUnlock()
	if caches, ok := mgr.redisMgr[component]; ok {
		return caches
	}
	return mgr.redisMgr[DefaultComponent]
}

func (mgr *Manager) GetCache(component string, key interface{}) redis.Cmdable {
	caches := mgr.GetCaches(component)
	if caches == nil || len(caches) == 0 {
		panic("no cache configuration for component: " + component)
	}
	sum := hash.Sum(key)
	return caches[sum%int64(len(caches))]
}
