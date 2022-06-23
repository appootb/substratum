package storage

import (
	"sync"

	"github.com/appootb/substratum/v2/storage"
)

func Init() {
	if storage.Implementor() == nil {
		storage.RegisterImplementor(&Manager{})
	}
}

type Manager struct {
	sync.Map
}

func (m *Manager) New(component string) {
	m.Store(component, &Storage{})
}

func (m *Manager) Get(component string) storage.Storage {
	s, ok := m.Load(component)
	if ok {
		return s.(storage.Storage)
	}
	return nil
}
