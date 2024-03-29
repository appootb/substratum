package storage

import (
	"sync"

	"github.com/appootb/substratum/v2/configure"
	"github.com/appootb/substratum/v2/storage"
)

func Init() {
	if storage.Implementor() == nil {
		storage.RegisterImplementor(&Manager{})
	}
	//
	storage.RegisterCommonDialectImplementor("", &emptyDialect{})
}

type Manager struct {
	sync.Map
}

func (m *Manager) New(component string) {
	m.Store(component, &Storage{
		common: make(map[configure.Schema]interface{}),
	})
}

func (m *Manager) Get(component string) storage.Storage {
	s, ok := m.Load(component)
	if ok {
		return s.(storage.Storage)
	}
	return nil
}
